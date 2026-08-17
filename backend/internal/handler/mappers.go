package handler

import (
	"github.com/supplychain/supplychain-api/internal/dto"
	"github.com/supplychain/supplychain-api/internal/model"
)

func toUserResponse(u model.User) dto.UserResponse {
	return dto.UserResponse{ID: u.ID, Username: u.Username, Role: string(u.Role), CreatedAt: u.CreatedAt}
}

func toSupplierResponse(s model.Supplier) dto.SupplierResponse {
	return dto.SupplierResponse{
		ID:            s.ID,
		Name:          s.Name,
		ContactPerson: s.ContactPerson,
		Phone:         s.Phone,
		Email:         s.Email,
		Address:       s.Address,
		Categories:    []string(s.Categories),
		Rating:        s.Rating,
		Status:        string(s.Status),
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}

func toInventoryResponse(item model.InventoryItem) dto.InventoryResponse {
	supplierName := ""
	if item.Supplier != nil {
		supplierName = item.Supplier.Name
	}
	return dto.InventoryResponse{
		ID:              item.ID,
		Name:            item.Name,
		Category:        string(item.Category),
		SupplierID:      item.SupplierID,
		SupplierName:    supplierName,
		BatchNo:         item.BatchNo,
		Quantity:        item.Quantity,
		Unit:            item.Unit,
		MinThreshold:    item.MinThreshold,
		ExpiryDate:      item.ExpiryDate.Format("2006-01-02"),
		StorageLocation: item.StorageLocation,
		Status:          string(item.Status),
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
	}
}

func toPurchaseOrderResponse(order model.PurchaseOrder) dto.PurchaseOrderResponse {
	supplierName := ""
	if order.Supplier != nil {
		supplierName = order.Supplier.Name
	}
	creatorName := ""
	if order.Creator != nil {
		creatorName = order.Creator.Username
	}
	approverName := ""
	if order.Approver != nil {
		approverName = order.Approver.Username
	}
	items := make([]dto.PurchaseOrderItemResponse, 0, len(order.Items))
	for _, it := range order.Items {
		name := ""
		if it.InventoryItem != nil {
			name = it.InventoryItem.Name
		}
		items = append(items, dto.PurchaseOrderItemResponse{
			ID:              it.ID,
			InventoryItemID: it.InventoryItemID,
			ItemName:        name,
			Quantity:        it.Quantity,
			UnitPrice:       it.UnitPrice,
			Subtotal:        it.Subtotal,
		})
	}
	return dto.PurchaseOrderResponse{
		ID:           order.ID,
		OrderNo:      order.OrderNo,
		SupplierID:   order.SupplierID,
		SupplierName: supplierName,
		Status:       string(order.Status),
		TotalAmount:  order.TotalAmount,
		CreatorID:    order.CreatorID,
		CreatorName:  creatorName,
		ApproverID:   order.ApproverID,
		ApproverName: approverName,
		ApprovedAt:   order.ApprovedAt,
		CompletedAt:  order.CompletedAt,
		Notes:        order.Notes,
		Items:        items,
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
	}
}

func toOperationLogResponse(l model.OperationLog) dto.OperationLogResponse {
	username := ""
	if l.User != nil {
		username = l.User.Username
	}
	return dto.OperationLogResponse{
		ID:         l.ID,
		UserID:     l.UserID,
		Username:   username,
		Action:     string(l.Action),
		TargetType: l.TargetType,
		TargetID:   l.TargetID,
		Detail:     l.Detail,
		CreatedAt:  l.CreatedAt,
	}
}
