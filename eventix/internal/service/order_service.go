package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"

	"eventix/internal/entity"
	"eventix/internal/repository"
	"eventix/pkg/worker"
)

var (
	ErrInsufficientTickets = errors.New("insufficient tickets available")
	ErrOrderNotFound       = errors.New("order not found")
	ErrOrderAlreadyPaid    = errors.New("order already paid")
	ErrOrderCancelled      = errors.New("order has been cancelled")
	ErrUnauthorized        = errors.New("unauthorized to access this order")
)

type OrderService interface {
	BookTickets(userID uint, eventID uint, qty int) (*entity.Order, error)
	ProcessPayment(userID uint, orderID uint) (*entity.Order, error)
	CancelOrder(userID uint, orderID uint) error
	GetUserOrders(userID uint) ([]entity.Order, error)
	GetOrderByID(userID uint, orderID uint) (*entity.Order, error)
}

type orderService struct {
	orderRepo    repository.OrderRepository
	eventRepo    repository.EventRepository
	ticketRepo   repository.TicketRepository
	emailChan    chan<- worker.EmailJob
	bookingMutex sync.Mutex
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	eventRepo repository.EventRepository,
	ticketRepo repository.TicketRepository,
	emailChan chan<- worker.EmailJob,
) OrderService {
	return &orderService{
		orderRepo:  orderRepo,
		eventRepo:  eventRepo,
		ticketRepo: ticketRepo,
		emailChan:  emailChan,
	}
}

// BookTickets handles ticket booking with mutex lock and database transaction
// to prevent overselling when multiple users book simultaneously
func (s *orderService) BookTickets(userID uint, eventID uint, qty int) (*entity.Order, error) {
	// CRITICAL SECTION: Lock to prevent race conditions
	s.bookingMutex.Lock()
	defer s.bookingMutex.Unlock()

	// Step 1: Check event exists and has enough tickets
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, ErrEventNotFound
	}

	if event.AvailableTickets < qty {
		return nil, ErrInsufficientTickets
	}

	// Step 2: Start database transaction
	db := s.orderRepo.GetDB()
	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Step 3: Decrement available tickets (within transaction)
	if err := s.eventRepo.DecrementAvailableTickets(tx, eventID, qty); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Double-check the affected rows to ensure tickets were actually decremented
	var updatedEvent entity.Event
	if err := tx.First(&updatedEvent, eventID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Step 4: Create order record
	order := &entity.Order{
		UserID:      userID,
		EventID:     eventID,
		Quantity:    qty,
		TotalAmount: event.Price * float64(qty),
		Status:      entity.OrderStatusPending,
	}

	if err := s.orderRepo.Save(tx, order); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Step 5: Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	order.Event = *event
	return order, nil
}

// ProcessPayment handles payment processing, ticket generation, and async email notification
func (s *orderService) ProcessPayment(userID uint, orderID uint) (*entity.Order, error) {
	// Step 1: Get and validate order
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	if order.UserID != userID {
		return nil, ErrUnauthorized
	}

	if order.Status == entity.OrderStatusPaid {
		return nil, ErrOrderAlreadyPaid
	}

	if order.Status == entity.OrderStatusCancelled {
		return nil, ErrOrderCancelled
	}

	// Step 2: Start transaction for payment processing
	db := s.orderRepo.GetDB()
	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Step 3: Update order status to PAID
	if err := s.orderRepo.UpdateStatus(tx, orderID, entity.OrderStatusPaid); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Step 4: Generate tickets
	tickets := make([]entity.Ticket, order.Quantity)
	for i := 0; i < order.Quantity; i++ {
		ticketCode, err := generateTicketCode()
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		tickets[i] = entity.Ticket{
			OrderID:    orderID,
			EventID:    order.EventID,
			TicketCode: ticketCode,
			Status:     entity.TicketStatusValid,
		}
	}

	if err := s.ticketRepo.SaveBatch(tx, tickets); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Step 5: Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Step 6: Fetch updated order with tickets
	updatedOrder, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	// Step 7: Send async email notification (fire and forget)
	// This runs in background via goroutine worker
	go func() {
		s.emailChan <- worker.EmailJob{
			Order: *updatedOrder,
			Email: fmt.Sprintf("user_%d@example.com", userID),
		}
	}()

	return updatedOrder, nil
}

// CancelOrder cancels a pending order and restores available tickets
func (s *orderService) CancelOrder(userID uint, orderID uint) error {
	s.bookingMutex.Lock()
	defer s.bookingMutex.Unlock()

	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return ErrOrderNotFound
	}

	if order.UserID != userID {
		return ErrUnauthorized
	}

	if order.Status == entity.OrderStatusPaid {
		return errors.New("cannot cancel a paid order")
	}

	if order.Status == entity.OrderStatusCancelled {
		return errors.New("order already cancelled")
	}

	db := s.orderRepo.GetDB()
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update order status
	if err := s.orderRepo.UpdateStatus(tx, orderID, entity.OrderStatusCancelled); err != nil {
		tx.Rollback()
		return err
	}

	// Restore available tickets
	if err := s.eventRepo.IncrementAvailableTickets(tx, order.EventID, order.Quantity); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *orderService) GetUserOrders(userID uint) ([]entity.Order, error) {
	return s.orderRepo.FindByUserID(userID)
}

func (s *orderService) GetOrderByID(userID uint, orderID uint) (*entity.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	if order.UserID != userID {
		return nil, ErrUnauthorized
	}

	return order, nil
}

func generateTicketCode() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "TKT-" + hex.EncodeToString(bytes), nil
}
