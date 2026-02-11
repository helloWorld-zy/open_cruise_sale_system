/**
 * @jest-environment jsdom
 */

import { mount } from '@vue/test-utils'
import PaymentResultPage from '../../pages/payment/result.vue'

describe('Payment Result Page', () => {
  beforeEach(() => {
    global.uni = {
      request: jest.fn(),
      redirectTo: jest.fn(),
      switchTab: jest.fn()
    }
    
    global.getApp = jest.fn(() => ({
      globalData: { apiBaseUrl: 'https://api.example.com' }
    }))
  })

  describe('Success State', () => {
    it('should display success message for successful payment', async () => {
      global.getCurrentPages = jest.fn(() => [{
        options: { order_id: 'order-123', status: 'success' }
      }])
      
      uni.request.mockResolvedValueOnce({
        data: {
          data: {
            id: 'order-123',
            order_number: 'ORD202401010001',
            total_amount: 2950,
            status: 'paid'
          }
        }
      })
      
      const wrapper = mount(PaymentResultPage)
      await new Promise(resolve => setTimeout(resolve, 0))
      
      expect(wrapper.vm.status).toBe('success')
      expect(wrapper.find('.title').text()).toBe('支付成功！')
    })

    it('should load order details on mount', async () => {
      global.getCurrentPages = jest.fn(() => [{
        options: { order_id: 'order-123' }
      }])
      
      const orderData = {
        id: 'order-123',
        order_number: 'ORD202401010001',
        total_amount: 2950,
        status: 'paid'
      }
      
      uni.request.mockResolvedValueOnce({
        data: { data: orderData }
      })
      
      const wrapper = mount(PaymentResultPage)
      await new Promise(resolve => setTimeout(resolve, 0))
      
      expect(uni.request).toHaveBeenCalledWith({
        url: 'https://api.example.com/orders/order-123'
      })
      expect(wrapper.vm.order).toEqual(orderData)
    })

    it('should navigate to order detail', async () => {
      global.getCurrentPages = jest.fn(() => [{
        options: { order_id: 'order-123', status: 'success' }
      }])
      
      const wrapper = mount(PaymentResultPage)
      wrapper.vm.orderId = 'order-123'
      
      wrapper.vm.goToOrder()
      
      expect(uni.redirectTo).toHaveBeenCalledWith({
        url: '/pages/orders/detail?id=order-123'
      })
    })

    it('should navigate to home', async () => {
      global.getCurrentPages = jest.fn(() => [{
        options: { order_id: 'order-123', status: 'success' }
      }])
      
      const wrapper = mount(PaymentResultPage)
      wrapper.vm.goHome()
      
      expect(uni.switchTab).toHaveBeenCalledWith({
        url: '/pages/index/index'
      })
    })
  })

  describe('Failed State', () => {
    it('should display failed message for failed payment', async () => {
      global.getCurrentPages = jest.fn(() => [{
        options: { order_id: 'order-123', status: 'failed' }
      }])
      
      const wrapper = mount(PaymentResultPage)
      await new Promise(resolve => setTimeout(resolve, 0))
      
      expect(wrapper.vm.status).toBe('failed')
      expect(wrapper.find('.title').text()).toBe('支付失败')
    })

    it('should allow retry payment', async () => {
      global.getCurrentPages = jest.fn(() => [{
        options: { order_id: 'order-123', status: 'failed' }
      }])
      
      const wrapper = mount(PaymentResultPage)
      wrapper.vm.orderId = 'order-123'
      
      wrapper.vm.retryPayment()
      
      expect(uni.redirectTo).toHaveBeenCalledWith({
        url: '/pages/booking/index?order_id=order-123'
      })
    })
  })

  describe('Payment Status Checking', () => {
    it('should poll order status until confirmed', async () => {
      global.getCurrentPages = jest.fn(() => [{
        options: { order_id: 'order-123' }
      }])
      
      // First call returns pending
      uni.request.mockResolvedValueOnce({
        data: { data: { id: 'order-123', status: 'pending' } }
      })
      
      // Second call returns paid
      uni.request.mockResolvedValueOnce({
        data: { data: { id: 'order-123', status: 'paid' } }
      })
      
      const wrapper = mount(PaymentResultPage)
      
      // Initial check
      await new Promise(resolve => setTimeout(resolve, 0))
      
      // Wait for polling interval
      jest.advanceTimersByTime(2000)
      await new Promise(resolve => setTimeout(resolve, 0))
      
      expect(wrapper.vm.status).toBe('success')
    })

    it('should handle API errors gracefully', async () => {
      global.getCurrentPages = jest.fn(() => [{
        options: { order_id: 'order-123' }
      }])
      
      uni.request.mockRejectedValueOnce(new Error('Network error'))
      
      const wrapper = mount(PaymentResultPage)
      await new Promise(resolve => setTimeout(resolve, 0))
      
      expect(wrapper.vm.status).toBe('failed')
    })
  })

  describe('Order Loading', () => {
    it('should handle missing order_id', async () => {
      global.getCurrentPages = jest.fn(() => [{
        options: {}
      }])
      
      const wrapper = mount(PaymentResultPage)
      await new Promise(resolve => setTimeout(resolve, 0))
      
      expect(wrapper.vm.status).toBe('failed')
    })
  })
})
