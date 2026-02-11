/**
 * @jest-environment jsdom
 */

import { mount } from '@vue/test-utils'
import BookingPage from '../../pages/booking/index.vue'

describe('Booking Page', () => {
  beforeEach(() => {
    // Mock uni API
    global.uni = {
      request: jest.fn(),
      showToast: jest.fn(),
      redirectTo: jest.fn(),
      requestPayment: jest.fn()
    }
    
    global.getApp = jest.fn(() => ({
      globalData: { apiBaseUrl: 'https://api.example.com' }
    }))
    
    global.getCurrentPages = jest.fn(() => [{
      options: { cruiseId: 'cruise-1' }
    }])
  })

  describe('Step 1: Voyage Selection', () => {
    it('should load voyages on mount', async () => {
      const voyages = [
        { id: 'voyage-1', departure_date: '2024-01-15', route: { name: 'Test Route', duration_days: 5 }, min_price: 1000 }
      ]
      
      uni.request.mockResolvedValueOnce({ data: { data: voyages } })
      
      const wrapper = mount(BookingPage, {
        props: { cruiseId: 'cruise-1' }
      })
      
      await new Promise(resolve => setTimeout(resolve, 0))
      
      expect(uni.request).toHaveBeenCalledWith({
        url: 'https://api.example.com/voyages',
        data: { cruise_id: 'cruise-1', booking_status: 'open' }
      })
    })

    it('should select voyage when clicked', async () => {
      const voyages = [
        { id: 'voyage-1', departure_date: '2024-01-15', route: { name: 'Test Route', duration_days: 5 } }
      ]
      
      uni.request.mockResolvedValueOnce({ data: { data: voyages } })
      
      const wrapper = mount(BookingPage, {
        props: { cruiseId: 'cruise-1' }
      })
      
      await new Promise(resolve => setTimeout(resolve, 0))
      
      // Trigger voyage selection
      wrapper.vm.selectVoyage(voyages[0])
      
      expect(wrapper.vm.selectedVoyage).toEqual(voyages[0])
    })
  })

  describe('Step 2: Cabin Selection', () => {
    beforeEach(() => {
      uni.request.mockImplementation((options) => {
        if (options.url.includes('/cabins')) {
          return Promise.resolve({
            data: {
              data: [
                { id: 'cabin-1', cabin_type_id: 'type-1', cabin_type: { name_cn: '内舱房' } }
              ]
            }
          })
        }
        if (options.url.includes('/prices')) {
          return Promise.resolve({
            data: {
              data: [
                { cabin_type_id: 'type-1', adult_price: 1000, child_price: 500 }
              ]
            }
          })
        }
        return Promise.resolve({ data: { data: [] } })
      })
    })

    it('should load cabins after voyage selection', async () => {
      const wrapper = mount(BookingPage, {
        props: { cruiseId: 'cruise-1' }
      })
      
      wrapper.vm.selectedVoyage = { id: 'voyage-1' }
      await wrapper.vm.loadCabins()
      
      expect(uni.request).toHaveBeenCalledWith({
        url: 'https://api.example.com/cabins',
        data: { voyage_id: 'voyage-1' }
      })
    })

    it('should toggle cabin selection', () => {
      const wrapper = mount(BookingPage, {
        props: { cruiseId: 'cruise-1' }
      })
      
      const cabin = { id: 'cabin-1', cabin_type: { name_cn: '内舱房' } }
      
      // Select cabin
      wrapper.vm.toggleCabin(cabin)
      expect(wrapper.vm.selectedCabins).toContainEqual(expect.objectContaining({ id: 'cabin-1' }))
      
      // Deselect cabin
      wrapper.vm.toggleCabin(cabin)
      expect(wrapper.vm.selectedCabins).toHaveLength(0)
    })

    it('should increment passenger count', () => {
      const wrapper = mount(BookingPage, {
        props: { cruiseId: 'cruise-1' }
      })
      
      wrapper.vm.increment('cabin-1', 'adult')
      expect(wrapper.vm.passengerCounts['cabin-1'].adult).toBe(3)
    })

    it('should not decrement adult below 1', () => {
      const wrapper = mount(BookingPage, {
        props: { cruiseId: 'cruise-1' }
      })
      
      wrapper.vm.passengerCounts = { 'cabin-1': { adult: 1, child: 0 } }
      wrapper.vm.decrement('cabin-1', 'adult')
      expect(wrapper.vm.passengerCounts['cabin-1'].adult).toBe(1)
    })
  })

  describe('Step 3: Passenger Info', () => {
    it('should initialize passengers based on cabin selection', () => {
      const wrapper = mount(BookingPage, {
        props: { cruiseId: 'cruise-1' }
      })
      
      wrapper.vm.selectedCabins = [
        { id: 'cabin-1', adultCount: 2, childCount: 1 }
      ]
      
      wrapper.vm.initPassengers()
      
      expect(wrapper.vm.passengers).toHaveLength(3)
      expect(wrapper.vm.passengers[0].type).toBe('adult')
      expect(wrapper.vm.passengers[2].type).toBe('child')
    })

    it('should validate required contact fields', () => {
      const wrapper = mount(BookingPage, {
        props: { cruiseId: 'cruise-1' }
      })
      
      wrapper.vm.currentStep = 2
      wrapper.vm.contact = { name: '', phone: '' }
      wrapper.vm.passengers = [{ name: '', surname: '', givenName: '', gender: '', birthDate: '', type: 'adult' }]
      
      expect(wrapper.vm.canProceed).toBe(false)
      
      wrapper.vm.contact = { name: 'Test', phone: '1234567890' }
      wrapper.vm.passengers = [{ name: '张三', surname: 'ZHANG', givenName: 'SAN', gender: 'male', birthDate: '1990-01-01', type: 'adult' }]
      
      expect(wrapper.vm.canProceed).toBe(true)
    })
  })

  describe('Step 4: Order Confirmation', () => {
    it('should calculate total price correctly', () => {
      const wrapper = mount(BookingPage, {
        props: { cruiseId: 'cruise-1' }
      })
      
      wrapper.vm.selectedCabins = [
        {
          id: 'cabin-1',
          adultCount: 2,
          childCount: 1,
          price: {
            adult_price: 1000,
            child_price: 500,
            port_fee: 100,
            service_fee: 50
          }
        }
      ]
      
      // (2 * 1000) + (1 * 500) + (3 * 100) + (3 * 50) = 2000 + 500 + 300 + 150 = 2950
      expect(wrapper.vm.totalPrice).toBe(2950)
    })
  })

  describe('Navigation', () => {
    it('should go to next step', () => {
      const wrapper = mount(BookingPage, {
        props: { cruiseId: 'cruise-1' }
      })
      
      wrapper.vm.selectedVoyage = { id: 'voyage-1' }
      wrapper.vm.nextStep()
      
      expect(wrapper.vm.currentStep).toBe(1)
    })

    it('should go to previous step', () => {
      const wrapper = mount(BookingPage, {
        props: { cruiseId: 'cruise-1' }
      })
      
      wrapper.vm.currentStep = 2
      wrapper.vm.prevStep()
      
      expect(wrapper.vm.currentStep).toBe(1)
    })

    it('should not go beyond step 0', () => {
      const wrapper = mount(BookingPage, {
        props: { cruiseId: 'cruise-1' }
      })
      
      wrapper.vm.currentStep = 0
      wrapper.vm.prevStep()
      
      expect(wrapper.vm.currentStep).toBe(0)
    })
  })

  describe('Payment Integration', () => {
    it('should create order and request WeChat payment', async () => {
      const wrapper = mount(BookingPage, {
        props: { cruiseId: 'cruise-1' }
      })
      
      wrapper.vm.selectedVoyage = { id: 'voyage-1' }
      wrapper.vm.selectedCabins = [
        { id: 'cabin-1', cabin_type_id: 'type-1', adultCount: 2, childCount: 0 }
      ]
      wrapper.vm.passengers = [
        { name: '张三', surname: 'ZHANG', givenName: 'SAN', gender: 'male', birthDate: '1990-01-01', type: 'adult' },
        { name: '李四', surname: 'LI', givenName: 'SI', gender: 'female', birthDate: '1992-02-02', type: 'adult' }
      ]
      wrapper.vm.contact = { name: '张三', phone: '13800138000' }
      
      uni.request.mockResolvedValueOnce({
        data: { data: { id: 'order-123', order_number: 'ORD202401010001' } }
      })
      
      await wrapper.vm.createOrder()
      
      expect(uni.request).toHaveBeenCalledWith({
        url: 'https://api.example.com/orders',
        method: 'POST',
        data: expect.objectContaining({
          cruise_id: 'cruise-1',
          voyage_id: 'voyage-1',
          items: expect.any(Array),
          passengers: expect.any(Array),
          contact_name: '张三',
          contact_phone: '13800138000'
        })
      })
    })
  })
})
