# Data Model: CruiseBooking Platform

**Version**: 1.0.0  
**Database**: PostgreSQL 17.x  
**ORM**: GORM v2.x  

---

## Entity Overview

Total: 18 core entities organized into 5 domains:

1. **Cruise Domain**: Companies, Ships, Facilities, Cabin Types
2. **Booking Domain**: Routes, Voyages, Cabins, Inventory, Prices
3. **Order Domain**: Orders, Order Items, Passengers, Payments, Refunds
4. **User Domain**: Users, Frequent Passengers, Staffs, Roles
5. **Content Domain**: Notifications, Travelogues

---

## Core Entities

### 1. cruise_companies

Cruise line operators (e.g., Royal Caribbean, Carnival)

```go
type CruiseCompany struct {
    ID          uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
    Name        string    `gorm:"not null;unique"`
    NameEN      string
    LogoURL     string
    Website     string
    Description string
    CreatedAt   time.Time `gorm:"autoCreateTime"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}
```

### 2. cruises

Individual cruise ships

```go
type Cruise struct {
    ID                uuid.UUID `gorm:"primaryKey"`
    CompanyID         uuid.UUID `gorm:"index"`
    Company           CruiseCompany
    NameCN            string `gorm:"not null"`
    NameEN            string
    Code              string `gorm:"not null;unique"`
    GrossTonnage      int
    PassengerCapacity int
    CrewCount         int
    BuiltYear         int
    RenovatedYear     int
    LengthMeters      float64
    WidthMeters       float64
    DeckCount         int
    CoverImages       datatypes.JSON
    Status            string `gorm:"default:active"` // active/inactive/maintenance
    SortWeight        int    `gorm:"default:0"`
    CreatedAt         time.Time
    UpdatedAt         time.Time
    DeletedAt         gorm.DeletedAt `gorm:"index"`
    
    // Relationships
    CabinTypes []CabinType
    Facilities []Facility
    Routes     []Route
}
```

### 3. cabin_types

Cabin categories (Inside, Oceanview, Balcony, Suite)

```go
type CabinType struct {
    ID              uuid.UUID `gorm:"primaryKey"`
    CruiseID        uuid.UUID `gorm:"index"`
    Cruise          Cruise
    Name            string `gorm:"not null"`
    Code            string `gorm:"not null"` // INSIDE, OCEANVIEW, etc.
    MinAreaSqm      float64
    MaxAreaSqm      float64
    StandardGuests  int
    MaxGuests       int
    BedTypes        string
    FeatureTags     datatypes.JSON
    Description     string
    Images          datatypes.JSON
    FloorPlanURL    string
    Amenities       datatypes.JSON
    SortWeight      int `gorm:"default:0"`
    Status          string `gorm:"default:active"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
    DeletedAt       gorm.DeletedAt `gorm:"index"`
}
```

### 4. facilities + facility_categories

Ship amenities

```go
type FacilityCategory struct {
    ID         uuid.UUID `gorm:"primaryKey"`
    CruiseID   uuid.UUID `gorm:"index"`
    Name       string `gorm:"not null"`
    Icon       string
    SortWeight int `gorm:"default:0"`
    CreatedAt  time.Time
    UpdatedAt  time.Time
    DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type Facility struct {
    ID           uuid.UUID `gorm:"primaryKey"`
    CruiseID     uuid.UUID `gorm:"index"`
    CategoryID   uuid.UUID `gorm:"index"`
    Category     FacilityCategory
    Name         string `gorm:"not null"`
    DeckNumber   int
    OpenTime     string
    IsFree       bool `gorm:"default:true"`
    Price        float64
    Description  string
    Images       datatypes.JSON
    SuitableTags datatypes.JSON
    SortWeight   int `gorm:"default:0"`
    Status       string `gorm:"default:visible"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    DeletedAt    gorm.DeletedAt `gorm:"index"`
}
```

### 5. routes

Cruise itineraries

```go
type Route struct {
    ID            uuid.UUID `gorm:"primaryKey"`
    CruiseID      uuid.UUID `gorm:"index"`
    Cruise        Cruise
    Name          string `gorm:"not null"`
    DeparturePort string
    ArrivalPort   string
    ViaPorts      datatypes.JSON
    DurationDays  int
    Description   string
    Status        string `gorm:"default:active"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
    DeletedAt     gorm.DeletedAt `gorm:"index"`
    
    Voyages []Voyage
}
```

### 6. voyages

Specific sailing instances

```go
type Voyage struct {
    ID             uuid.UUID `gorm:"primaryKey"`
    RouteID        uuid.UUID `gorm:"index"`
    Route          Route
    DepartureDate  time.Time `gorm:"not null"`
    ReturnDate     time.Time
    Status         string `gorm:"default:open"` // open/closed/cancelled
    CreatedAt      time.Time
    UpdatedAt      time.Time
    DeletedAt      gorm.DeletedAt `gorm:"index"`
}
```

### 7. cabins

Individual bookable units

```go
type Cabin struct {
    ID           uuid.UUID `gorm:"primaryKey"`
    CabinTypeID  uuid.UUID `gorm:"index"`
    CabinType    CabinType
    VoyageID     uuid.UUID `gorm:"index"`
    Voyage       Voyage
    CabinNumber  string `gorm:"not null"`
    Grade        string
    BedType      string
    AreaSqm      float64
    DeckNumber   int
    Position     string // front/middle/rear
    Orientation  string // port/starboard
    HasWindow    bool `gorm:"default:false"`
    HasBalcony   bool `gorm:"default:false"`
    Amenities    datatypes.JSON
    Status       string `gorm:"default:available"` // available/locked/sold/offline
    CreatedAt    time.Time
    UpdatedAt    time.Time
    DeletedAt    gorm.DeletedAt `gorm:"index"`
}
```

### 8. cabin_inventory

Real-time stock tracking

```go
type CabinInventory struct {
    ID              uuid.UUID `gorm:"primaryKey"`
    CabinTypeID     uuid.UUID `gorm:"uniqueIndex:idx_cabin_type_voyage"`
    CabinType       CabinType
    VoyageID        uuid.UUID `gorm:"uniqueIndex:idx_cabin_type_voyage"`
    Voyage          Voyage
    TotalQuantity   int `gorm:"not null"`
    SoldQuantity    int `gorm:"default:0"`
    LockedQuantity  int `gorm:"default:0"`
    AlertThreshold  int `gorm:"default:5"`
    Version         int `gorm:"default:1"` // Optimistic locking
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

// Available = Total - Sold - Locked
```

### 9. cabin_prices

Date-specific pricing

```go
type CabinPrice struct {
    ID                  uuid.UUID `gorm:"primaryKey"`
    CabinTypeID         uuid.UUID `gorm:"index"`
    VoyageID            uuid.UUID `gorm:"index"`
    Date                time.Time `gorm:"not null"`
    BasePrice           float64 `gorm:"not null"`
    ChildPrice          float64
    SingleSupplement    float64
    HolidayMarkup       float64 `gorm:"default:0"`
    EarlyBirdDiscount   float64 `gorm:"default:0"`
    ConsecutiveDiscount float64 `gorm:"default:0"`
    Currency            string  `gorm:"default:CNY"`
    CreatedAt           time.Time
    UpdatedAt           time.Time
}
```

### 10. orders

Customer orders

```go
type Order struct {
    ID              uuid.UUID `gorm:"primaryKey"`
    OrderNo         string    `gorm:"not null;unique"` // CB-YYYYMMDD-XXXXX
    UserID          uuid.UUID `gorm:"index"`
    User            User
    Status          string `gorm:"default:created"`
    TotalAmount     float64 `gorm:"not null"`
    DiscountAmount  float64 `gorm:"default:0"`
    PaidAmount      float64 `gorm:"default:0"`
    Currency        string `gorm:"default:CNY"`
    DepartureDate   time.Time
    Notes           string
    PaidAt          *time.Time
    ConfirmedAt     *time.Time
    CancelledAt     *time.Time
    CompletedAt     *time.Time
    ExpireAt        time.Time // 15 min from creation
    CreatedAt       time.Time
    UpdatedAt       time.Time
    DeletedAt       gorm.DeletedAt `gorm:"index"`
    
    Items      []OrderItem
    Passengers []Passenger
    Payments   []Payment
}

// Status flow:
// created → pending_payment → paid → confirmed → pending_departure → departed → completed
//    ↓           ↓              ↓
// cancelled   timeout      refund_requested → refund_processing → refunded
```

### 11. order_items

Line items in orders

```go
type OrderItem struct {
    ID          uuid.UUID `gorm:"primaryKey"`
    OrderID     uuid.UUID `gorm:"index"`
    Order       Order
    CabinID     uuid.UUID
    Cabin       Cabin
    VoyageID    uuid.UUID
    Voyage      Voyage
    Quantity    int `gorm:"not null"`
    UnitPrice   float64 `gorm:"not null"`
    TotalPrice  float64 `gorm:"not null"`
    GuestCount  int
    CreatedAt   time.Time
}
```

### 12. passengers

Guest information

```go
type Passenger struct {
    ID              uuid.UUID `gorm:"primaryKey"`
    OrderID         uuid.UUID `gorm:"index"`
    Order           Order
    OrderItemID     uuid.UUID `gorm:"index"`
    NameCN          string `gorm:"not null"`
    NameEN          string
    IDType          string `gorm:"not null"` // id_card/passport/other
    IDNumber        string `gorm:"not null"` // Encrypted
    Phone           string
    Email           string
    EmergencyName   string
    EmergencyPhone  string
    SpecialRequests string
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### 13. payments

Transaction records

```go
type Payment struct {
    ID            uuid.UUID `gorm:"primaryKey"`
    OrderID       uuid.UUID `gorm:"index"`
    Order         Order
    PaymentNo     string `gorm:"not null;unique"`
    Channel       string `gorm:"not null"` // wechat/alipay/unionpay
    Amount        float64 `gorm:"not null"`
    Status        string `gorm:"default:pending"`
    ThirdPartyID  string
    PaidAt        *time.Time
    CallbackData  datatypes.JSON
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

### 14. users

Customer accounts

```go
type User struct {
    ID             uuid.UUID `gorm:"primaryKey"`
    UserNo         string    `gorm:"not null;unique"`
    WechatOpenID   string    `gorm:"index"`
    WechatUnionID  string
    Phone          string    `gorm:"unique"`
    Email          string
    Nickname       string
    AvatarURL      string
    MemberLevel    int `gorm:"default:1"`
    PointsBalance  int `gorm:"default:0"`
    Status         string `gorm:"default:active"`
    CreatedAt      time.Time
    UpdatedAt      time.Time
    DeletedAt      gorm.DeletedAt `gorm:"index"`
    
    Orders            []Order
    FrequentPassengers []FrequentPassenger
}
```

### 15. frequent_passengers

Saved guest profiles

```go
type FrequentPassenger struct {
    ID           uuid.UUID `gorm:"primaryKey"`
    UserID       uuid.UUID `gorm:"index"`
    User         User
    Name         string `gorm:"not null"`
    NameEN       string
    IDType       string
    IDNumber     string // Encrypted
    Phone        string
    Email        string
    Relationship string // self/spouse/child/parent
    CreatedAt    time.Time
    UpdatedAt    time.Time
    DeletedAt    gorm.DeletedAt `gorm:"index"`
}
```

### 16. staffs + roles

Admin RBAC

```go
type Role struct {
    ID          uuid.UUID `gorm:"primaryKey"`
    Name        string `gorm:"not null"`
    Code        string `gorm:"not null;unique"`
    Permissions datatypes.JSON
    Description string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Staff struct {
    ID           uuid.UUID `gorm:"primaryKey"`
    RoleID       uuid.UUID `gorm:"index"`
    Role         Role
    Username     string `gorm:"not null;unique"`
    PasswordHash string `gorm:"not null"`
    Name         string
    Email        string
    Phone        string
    LastLoginAt  *time.Time
    Status       string `gorm:"default:active"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    DeletedAt    gorm.DeletedAt `gorm:"index"`
}
```

### 17. notifications

User messages

```go
type Notification struct {
    ID        uuid.UUID `gorm:"primaryKey"`
    UserID    uuid.UUID `gorm:"index"`
    User      User
    Type      string `gorm:"not null"` // order/promotion/system
    Channel   string `gorm:"not null"` // wechat/sms/in_app
    Title     string
    Content   string
    IsRead    bool `gorm:"default:false"`
    ReadAt    *time.Time
    CreatedAt time.Time
}
```

### 18. travelogues

Community content (V2.0)

```go
type Travelogue struct {
    ID          uuid.UUID `gorm:"primaryKey"`
    UserID      uuid.UUID `gorm:"index"`
    User        User
    Title       string `gorm:"not null"`
    Content     string
    Images      datatypes.JSON
    RouteID     *uuid.UUID
    LikesCount  int `gorm:"default:0"`
    SavesCount  int `gorm:"default:0"`
    CommentsCount int `gorm:"default:0"`
    Status      string `gorm:"default:pending"` // pending/approved/rejected
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}
```

---

## Indexes Summary

### Primary Keys
All entities use `uuid` as primary key with `gen_random_uuid()` default

### Unique Constraints
- cruise_companies.name
- cruises.code
- cabin_types (cruise_id + code)
- voyages (route_id + departure_date)
- cabins (voyage_id + cabin_number)
- cabin_inventory (cabin_type_id + voyage_id)
- orders.order_no
- payments.payment_no
- roles.code
- staffs.username
- users (user_no, phone, wechat_openid)

### Performance Indexes
- cruises (company_id, status, sort_weight)
- cabin_types (cruise_id, status, sort_weight)
- voyages (route_id, departure_date, status)
- orders (user_id, status, created_at)
- notifications (user_id, is_read, created_at)

---

## Soft Delete Strategy

All business entities implement soft delete using GORM's `DeletedAt`:

```go
DeletedAt gorm.DeletedAt `gorm:"index"`
```

Benefits:
- Data audit trail
- Recover accidentally deleted data
- Maintain referential integrity
- Compliance requirements

---

## Data Integrity Rules

1. **Inventory Non-Negative**: `available = total - sold - locked >= 0`
2. **Order Amount Consistency**: `total = sum(items) - discounts`
3. **Payment Amount Validation**: `payment.amount == order.total_amount`
4. **Unique Order Number**: Format `CB-YYYYMMDD-XXXXX`
5. **Unique Payment Number**: Format `PAY-YYYYMMDD-XXXXX`

---

**End of Data Model**
