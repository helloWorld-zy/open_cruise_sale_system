// Auth Types
export interface User {
  id: string
  userNo: string
  phone?: string
  email?: string
  nickname?: string
  avatarUrl?: string
  memberLevel: number
  pointsBalance: number
  status: 'active' | 'suspended'
  createdAt: string
  updatedAt: string
}

export interface Staff {
  id: string
  username: string
  name: string
  email?: string
  phone?: string
  role: Role
  lastLoginAt?: string
  status: 'active' | 'inactive'
}

export interface Role {
  id: string
  name: string
  code: string
  permissions: Permission[]
  description?: string
}

export interface Permission {
  resource: string
  action: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  refreshToken: string
  user: User | Staff
}

// Cruise Types
export interface Cruise {
  id: string
  companyId: string
  nameCn: string
  nameEn?: string
  code: string
  grossTonnage?: number
  passengerCapacity?: number
  crewCount?: number
  builtYear?: number
  renovatedYear?: number
  lengthMeters?: number
  widthMeters?: number
  deckCount?: number
  coverImages: string[]
  status: 'active' | 'inactive' | 'maintenance'
  sortWeight: number
  createdAt: string
  updatedAt: string
}

export interface CruiseCompany {
  id: string
  name: string
  nameEn?: string
  logoUrl?: string
  website?: string
  description?: string
}

export interface CabinType {
  id: string
  cruiseId: string
  name: string
  code: string
  minAreaSqm?: number
  maxAreaSqm?: number
  standardGuests: number
  maxGuests: number
  bedTypes?: string
  featureTags: string[]
  description?: string
  images: string[]
  floorPlanUrl?: string
  amenities: string[]
  sortWeight: number
  status: 'active' | 'inactive'
}

export interface Facility {
  id: string
  cruiseId: string
  categoryId: string
  name: string
  deckNumber?: number
  openTime?: string
  isFree: boolean
  price?: number
  description?: string
  images: string[]
  suitableTags: string[]
  sortWeight: number
  status: 'visible' | 'hidden'
}

export interface FacilityCategory {
  id: string
  cruiseId: string
  name: string
  icon?: string
  sortWeight: number
}

// Route & Voyage Types
export interface Route {
  id: string
  cruiseId: string
  name: string
  departurePort?: string
  arrivalPort?: string
  viaPorts: string[]
  durationDays: number
  description?: string
  status: 'active' | 'inactive'
}

export interface Voyage {
  id: string
  routeId: string
  departureDate: string
  returnDate?: string
  status: 'open' | 'closed' | 'cancelled'
}

// Cabin & Inventory Types
export interface Cabin {
  id: string
  cabinTypeId: string
  voyageId: string
  cabinNumber: string
  grade?: string
  bedType?: string
  areaSqm?: number
  deckNumber?: number
  position?: 'front' | 'middle' | 'rear'
  orientation?: 'port' | 'starboard'
  hasWindow: boolean
  hasBalcony: boolean
  amenities: string[]
  status: 'available' | 'locked' | 'sold' | 'offline'
}

export interface CabinInventory {
  id: string
  cabinTypeId: string
  voyageId: string
  totalQuantity: number
  soldQuantity: number
  lockedQuantity: number
  availableQuantity: number
  alertThreshold: number
}

export interface CabinPrice {
  id: string
  cabinTypeId: string
  voyageId: string
  date: string
  basePrice: number
  childPrice?: number
  singleSupplement?: number
  holidayMarkup: number
  earlyBirdDiscount: number
  consecutiveDiscount: number
  currency: string
}

// Order Types
export type OrderStatus = 
  | 'created'
  | 'pending_payment'
  | 'paid'
  | 'confirmed'
  | 'pending_departure'
  | 'departed'
  | 'completed'
  | 'cancelled'
  | 'refund_requested'
  | 'refund_processing'
  | 'refunded'

export interface Order {
  id: string
  orderNo: string
  userId: string
  status: OrderStatus
  totalAmount: number
  discountAmount: number
  paidAmount: number
  currency: string
  departureDate?: string
  notes?: string
  paidAt?: string
  confirmedAt?: string
  cancelledAt?: string
  completedAt?: string
  expireAt: string
  createdAt: string
  updatedAt: string
  items: OrderItem[]
  passengers: Passenger[]
  payments: Payment[]
}

export interface OrderItem {
  id: string
  orderId: string
  cabinId: string
  voyageId: string
  quantity: number
  unitPrice: number
  totalPrice: number
  guestCount?: number
  cabin?: Cabin
}

export interface Passenger {
  id: string
  orderId: string
  orderItemId?: string
  nameCn: string
  nameEn?: string
  idType: 'id_card' | 'passport' | 'other'
  idNumber: string
  phone?: string
  email?: string
  emergencyName?: string
  emergencyPhone?: string
  specialRequests?: string
}

// Payment Types
export type PaymentStatus = 'pending' | 'processing' | 'success' | 'failed' | 'refunded'
export type PaymentChannel = 'wechat' | 'alipay' | 'unionpay'

export interface Payment {
  id: string
  orderId: string
  paymentNo: string
  channel: PaymentChannel
  amount: number
  status: PaymentStatus
  thirdPartyId?: string
  paidAt?: string
  createdAt: string
}

// Notification Types
export interface Notification {
  id: string
  userId: string
  type: 'order' | 'promotion' | 'system'
  channel: 'wechat' | 'sms' | 'in_app'
  title: string
  content: string
  isRead: boolean
  readAt?: string
  createdAt: string
}

// API Response Types
export interface ApiResponse<T> {
  code: number
  message: string
  data?: T
  error?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  pagination: {
    page: number
    pageSize: number
    total: number
    pages: number
  }
}

// Frequent Passenger
export interface FrequentPassenger {
  id: string
  userId: string
  name: string
  nameEn?: string
  idType: string
  idNumber: string
  phone?: string
  email?: string
  relationship?: string
}
