import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import productService from '../../services/productService';

const PRODUCT_QUERY_KEYS = {
  ALL: ['products'],
  LISTS: () => [...PRODUCT_QUERY_KEYS.ALL, 'list'],
  LIST: (filters) => [...PRODUCT_QUERY_KEYS.LISTS(), { filters }],
  DETAILS: () => [...PRODUCT_QUERY_KEYS.ALL, 'detail'],
  DETAIL: (id) => [...PRODUCT_QUERY_KEYS.DETAILS(), id],
};

export const useProducts = (filters = {}) => {
  return useQuery({
    queryKey: PRODUCT_QUERY_KEYS.LIST(filters),
    queryFn: () => productService.getProducts(filters),
    keepPreviousData: true,
  });
};

export const useProduct = (id) => {
  return useQuery({
    queryKey: PRODUCT_QUERY_KEYS.DETAIL(id),
    queryFn: () => productService.getProduct(id),
    enabled: !!id,
  });
};

export const useCreateProduct = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (productData) => productService.createProduct(productData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PRODUCT_QUERY_KEYS.LISTS() });
    },
  });
};

export const useUpdateProduct = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, data }) => productService.updateProduct(id, data),
    onSuccess: (data, { id }) => {
      queryClient.invalidateQueries({ queryKey: PRODUCT_QUERY_KEYS.DETAIL(id) });
      queryClient.invalidateQueries({ queryKey: PRODUCT_QUERY_KEYS.LISTS() });
    },
  });
};

export const useDeleteProduct = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (id) => productService.deleteProduct(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PRODUCT_QUERY_KEYS.LISTS() });
    },
  });
};
