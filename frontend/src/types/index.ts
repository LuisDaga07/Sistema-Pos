export interface Product {
  id: string;
  restaurant_id: string;
  category_id?: string;
  name: string;
  description: string;
  price: number;
  image_url?: string;
  active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: string;
  restaurant_id: string;
  name: string;
  description: string;
  sort_order: number;
}

export interface Sale {
  id: string;
  restaurant_id: string;
  user_id: string;
  total: number;
  status: string;
  created_at: string;
  updated_at: string;
}
