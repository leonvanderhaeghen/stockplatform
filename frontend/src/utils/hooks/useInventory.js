import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import inventoryService from '../../services/inventoryService';

const INVENTORY_QUERY_KEYS = {
  ALL: ['inventory'],
  LISTS: () => [...INVENTORY_QUERY_KEYS.ALL, 'list'],
  LIST: (filters) => [...INVENTORY_QUERY_KEYS.LISTS(), { filters }],
  DETAILS: () => [...INVENTORY_QUERY_KEYS.ALL, 'detail'],
  DETAIL: (id) => [...INVENTORY_QUERY_KEYS.DETAILS(), id],
  HISTORY: (id) => [...INVENTORY_QUERY_KEYS.ALL, 'history', id],
  RESERVATIONS: () => [...INVENTORY_QUERY_KEYS.ALL, 'reservations'],
  MOVEMENTS: () => [...INVENTORY_QUERY_KEYS.ALL, 'movements'],
  LOW_STOCK: () => [...INVENTORY_QUERY_KEYS.ALL, 'low-stock'],
};

// Get all inventory items with optional filters
export const useInventoryItems = (filters = {}) => {
  return useQuery({
    queryKey: INVENTORY_QUERY_KEYS.LIST(filters),
    queryFn: async () => {
      const response = await inventoryService.getInventoryItems(filters);
      // Handle both response formats: { data: [...] } and direct array
      if (Array.isArray(response)) return response;
      if (response && Array.isArray(response.data)) return response.data;
      return [];
    },
    keepPreviousData: true,
  });
};

// Get a single inventory item by ID
export const useInventoryItem = (id) => {
  return useQuery({
    queryKey: INVENTORY_QUERY_KEYS.DETAIL(id),
    queryFn: () => inventoryService.getInventoryItem(id),
    enabled: !!id,
    retry: false, // Don't retry if inventory item doesn't exist
  });
};

// Get inventory history for an item
export const useInventoryHistory = (id) => {
  return useQuery({
    queryKey: INVENTORY_QUERY_KEYS.HISTORY(id),
    queryFn: () => inventoryService.getInventoryHistory(id),
    enabled: !!id,
  });
};

// Get inventory reservations
export const useInventoryReservations = (params = {}) => {
  return useQuery({
    queryKey: [...INVENTORY_QUERY_KEYS.RESERVATIONS(), params],
    queryFn: () => inventoryService.getReservations(params),
    keepPreviousData: true,
  });
};

// Get inventory movements
export const useInventoryMovements = (params = {}) => {
  return useQuery({
    queryKey: [...INVENTORY_QUERY_KEYS.MOVEMENTS(), params],
    queryFn: () => inventoryService.getInventoryMovements(params),
    keepPreviousData: true,
  });
};

// Get low stock items
export const useLowStockItems = (threshold = 10, location = '') => {
  return useQuery({
    queryKey: [...INVENTORY_QUERY_KEYS.LOW_STOCK(), { threshold, location }],
    queryFn: () => inventoryService.getLowStockItems(threshold, location),
  });
};

// Create inventory item mutation
export const useCreateInventoryItem = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (itemData) => inventoryService.createInventoryItem(itemData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.LISTS() });
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.LOW_STOCK() });
    },
  });
};

// Update inventory item mutation
export const useUpdateInventoryItem = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, data }) => inventoryService.updateInventoryItem(id, data),
    onSuccess: (data, { id }) => {
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.DETAIL(id) });
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.LISTS() });
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.LOW_STOCK() });
    },
  });
};

// Update inventory quantity mutation
export const useUpdateInventoryQuantity = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, change, reason }) => inventoryService.updateInventoryQuantity(id, change, reason),
    onSuccess: (data, { id }) => {
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.DETAIL(id) });
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.LISTS() });
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.HISTORY(id) });
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.MOVEMENTS() });
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.LOW_STOCK() });
    },
  });
};

// Delete inventory item mutation
export const useDeleteInventoryItem = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (id) => inventoryService.deleteInventoryItem(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.LISTS() });
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.LOW_STOCK() });
    },
  });
};

// Reserve inventory mutation
export const useReserveInventory = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (reservationData) => inventoryService.reserveInventory(reservationData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.LISTS() });
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.RESERVATIONS() });
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.MOVEMENTS() });
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.LOW_STOCK() });
    },
  });
};

// Release inventory mutation
export const useReleaseInventory = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (releaseData) => inventoryService.releaseInventory(releaseData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.LISTS() });
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.RESERVATIONS() });
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.MOVEMENTS() });
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.LOW_STOCK() });
    },
  });
};

// Bulk update inventory mutation
export const useBulkUpdateInventory = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (bulkData) => inventoryService.bulkUpdateInventory(bulkData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.LISTS() });
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.MOVEMENTS() });
      queryClient.invalidateQueries({ queryKey: INVENTORY_QUERY_KEYS.LOW_STOCK() });
    },
  });
};

const inventoryHooks = {
  useInventoryItems,
  useInventoryItem,
  useInventoryHistory,
  useInventoryReservations,
  useInventoryMovements,
  useLowStockItems,
  useCreateInventoryItem,
  useUpdateInventoryItem,
  useUpdateInventoryQuantity,
  useDeleteInventoryItem,
  useReserveInventory,
  useReleaseInventory,
  useBulkUpdateInventory,
};

export default inventoryHooks;
