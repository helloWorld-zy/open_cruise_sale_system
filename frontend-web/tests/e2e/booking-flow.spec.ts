import { test, expect } from '@playwright/test'

test.describe('Booking Flow', () => {
  test.describe('Step 1: Select Voyage', () => {
    test.beforeEach(async ({ page }) => {
      // Navigate to a cruise detail page with booking option
      await page.goto('/cruises/test-cruise-id')
      await page.waitForSelector('[data-testid="cruise-detail"]')
    })

    test('should display voyage selection step', async ({ page }) => {
      // Click book now button
      await page.click('text=立即预订')
      
      // Wait for booking wizard to load
      await page.waitForURL('/booking**')
      await page.waitForSelector('[data-testid="booking-wizard"]')
      
      // Check step indicator shows step 1
      await expect(page.locator('[data-testid="step-indicator"]')).toContainText('选择航次')
      
      // Check voyage list is displayed
      await expect(page.locator('[data-testid="voyage-list"]')).toBeVisible()
      await expect(page.locator('[data-testid="voyage-card"]').first()).toBeVisible()
    })

    test('should select a voyage and proceed to step 2', async ({ page }) => {
      await page.goto('/booking?cruiseId=test-cruise-id')
      await page.waitForSelector('[data-testid="voyage-list"]')
      
      // Click on first voyage card
      await page.click('[data-testid="voyage-card"]:first-child')
      
      // Check voyage is selected (has selected class or visual indicator)
      await expect(page.locator('[data-testid="voyage-card"].selected')).toBeVisible()
      
      // Click next button
      await page.click('button:has-text("下一步")')
      
      // Should navigate to step 2
      await page.waitForSelector('[data-testid="cabin-selection"]')
      await expect(page.locator('[data-testid="step-indicator"]')).toContainText('选择舱房')
    })

    test('should not proceed without selecting voyage', async ({ page }) => {
      await page.goto('/booking?cruiseId=test-cruise-id')
      
      // Try to click next without selecting
      const nextButton = page.locator('button:has-text("下一步")')
      await expect(nextButton).toBeDisabled()
      
      // Should still be on step 1
      await expect(page.locator('[data-testid="voyage-list"]')).toBeVisible()
    })

    test('should display voyage details correctly', async ({ page }) => {
      await page.goto('/booking?cruiseId=test-cruise-id')
      
      // Check voyage card contains required info
      const voyageCard = page.locator('[data-testid="voyage-card"]').first()
      await expect(voyageCard).toContainText('出发')
      await expect(voyageCard).toContainText('天航程')
      await expect(voyageCard.locator('.price')).toBeVisible()
    })
  })

  test.describe('Step 2: Select Cabin', () => {
    test.beforeEach(async ({ page }) => {
      // Navigate to booking with voyage pre-selected
      await page.goto('/booking?cruiseId=test-cruise-id&voyageId=test-voyage-id')
      await page.waitForSelector('[data-testid="cabin-selection"]')
    })

    test('should display cabin selection step', async ({ page }) => {
      // Check cabin list is displayed
      await expect(page.locator('[data-testid="cabin-list"]')).toBeVisible()
      await expect(page.locator('[data-testid="cabin-card"]').first()).toBeVisible()
      
      // Check passenger counters are present
      await expect(page.locator('[data-testid="adult-counter"]')).toBeVisible()
      await expect(page.locator('[data-testid="child-counter"]')).toBeVisible()
    })

    test('should increment and decrement passenger count', async ({ page }) => {
      // Get initial adult count (should be 2)
      const adultCount = page.locator('[data-testid="adult-counter"] .count')
      await expect(adultCount).toHaveText('2')
      
      // Click increment
      await page.click('[data-testid="adult-counter"] button:has-text("+")')
      await expect(adultCount).toHaveText('3')
      
      // Click decrement
      await page.click('[data-testid="adult-counter"] button:has-text("-")')
      await expect(adultCount).toHaveText('2')
      
      // Try to decrement below 1 (should not work)
      await page.click('[data-testid="adult-counter"] button:has-text("-")')
      await page.click('[data-testid="adult-counter"] button:has-text("-")')
      await expect(adultCount).toHaveText('1') // Should stop at 1
    })

    test('should select cabin and update summary', async ({ page }) => {
      // Click select cabin button
      await page.click('[data-testid="cabin-card"]:first-child button:has-text("选择此舱房")')
      
      // Check button changes to "已选择"
      await expect(page.locator('[data-testid="cabin-card"]:first-child button')).toContainText('已选择')
      
      // Check summary shows selected cabin
      await expect(page.locator('[data-testid="selected-summary"]')).toContainText('已选择舱房')
      await expect(page.locator('[data-testid="total-price"]')).toBeVisible()
    })

    test('should calculate price correctly', async ({ page }) => {
      // Select a cabin
      await page.click('[data-testid="cabin-card"]:first-child button:has-text("选择此舱房")')
      
      // Set 2 adults, 1 child
      await page.click('[data-testid="adult-counter"] button:has-text("+")')
      await page.click('[data-testid="child-counter"] button:has-text("+")')
      
      // Check price calculation
      const priceText = await page.locator('[data-testid="total-price"]').textContent()
      expect(priceText).toMatch(/¥[\d,]+/)
    })

    test('should proceed to step 3 after selecting cabin', async ({ page }) => {
      // Select cabin
      await page.click('[data-testid="cabin-card"]:first-child button:has-text("选择此舱房")')
      
      // Click next
      await page.click('button:has-text("下一步")')
      
      // Should navigate to step 3
      await page.waitForSelector('[data-testid="passenger-form"]')
      await expect(page.locator('[data-testid="step-indicator"]')).toContainText('填写乘客')
    })

    test('should allow going back to step 1', async ({ page }) => {
      // Click previous
      await page.click('button:has-text("上一步")')
      
      // Should be back on voyage selection
      await page.waitForSelector('[data-testid="voyage-list"]')
      await expect(page.locator('[data-testid="step-indicator"]')).toContainText('选择航次')
    })
  })

  test.describe('Step 3: Passenger Information', () => {
    test.beforeEach(async ({ page }) => {
      // Navigate to step 3 with voyage and cabin pre-selected
      await page.goto('/booking?cruiseId=test-cruise-id&voyageId=test-voyage-id&cabinId=test-cabin-id')
      await page.waitForSelector('[data-testid="passenger-form"]')
    })

    test('should display passenger form step', async ({ page }) => {
      // Check contact information section
      await expect(page.locator('[data-testid="contact-section"]')).toBeVisible()
      await expect(page.locator('input[name="contactName"]')).toBeVisible()
      await expect(page.locator('input[name="contactPhone"]')).toBeVisible()
      
      // Check passenger cards are displayed
      await expect(page.locator('[data-testid="passenger-card"]').first()).toBeVisible()
    })

    test('should validate required passenger fields', async ({ page }) => {
      // Try to proceed without filling required fields
      await page.click('button:has-text("下一步")')
      
      // Should show validation errors
      await expect(page.locator('text=请输入联系人姓名')).toBeVisible()
      await expect(page.locator('text=请输入手机号码')).toBeVisible()
    })

    test('should fill passenger information correctly', async ({ page }) => {
      // Fill contact info
      await page.fill('input[name="contactName"]', '张三')
      await page.fill('input[name="contactPhone"]', '13800138000')
      
      // Fill first passenger
      await page.fill('[data-testid="passenger-card"]:first-child input[name="surname"]', 'ZHANG')
      await page.fill('[data-testid="passenger-card"]:first-child input[name="givenName"]', 'SAN')
      await page.fill('[data-testid="passenger-card"]:first-child input[name="name"]', '张三')
      await page.selectOption('[data-testid="passenger-card"]:first-child select[name="gender"]', 'male')
      await page.fill('[data-testid="passenger-card"]:first-child input[name="birthDate"]', '1990-01-01')
      
      // Validation errors should be gone
      await expect(page.locator('text=请输入联系人姓名')).not.toBeVisible()
    })

    test('should proceed to step 4 after filling all required info', async ({ page }) => {
      // Fill contact info
      await page.fill('input[name="contactName"]', '张三')
      await page.fill('input[name="contactPhone"]', '13800138000')
      
      // Fill all passengers
      const passengerCards = page.locator('[data-testid="passenger-card"]')
      const count = await passengerCards.count()
      
      for (let i = 0; i < count; i++) {
        await passengerCards.nth(i).locator('input[name="surname"]').fill(`SURNAME${i}`)
        await passengerCards.nth(i).locator('input[name="givenName"]').fill(`GIVEN${i}`)
        await passengerCards.nth(i).locator('input[name="name"]').fill(`乘客${i}`)
        await passengerCards.nth(i).locator('select[name="gender"]').selectOption('male')
        await passengerCards.nth(i).locator('input[name="birthDate"]').fill('1990-01-01')
      }
      
      // Click next
      await page.click('button:has-text("下一步")')
      
      // Should navigate to step 4
      await page.waitForSelector('[data-testid="order-confirmation"]')
      await expect(page.locator('[data-testid="step-indicator"]')).toContainText('确认订单')
    })
  })

  test.describe('Step 4: Order Confirmation', () => {
    test.beforeEach(async ({ page }) => {
      // Navigate to step 4 with all info pre-filled
      await page.goto('/booking?cruiseId=test-cruise-id&step=4')
      await page.waitForSelector('[data-testid="order-confirmation"]')
    })

    test('should display order confirmation step', async ({ page }) => {
      // Check order summary sections
      await expect(page.locator('[data-testid="voyage-summary"]')).toBeVisible()
      await expect(page.locator('[data-testid="cabin-summary"]')).toBeVisible()
      await expect(page.locator('[data-testid="passenger-summary"]')).toBeVisible()
      await expect(page.locator('[data-testid="price-summary"]')).toBeVisible()
    })

    test('should display correct order details', async ({ page }) => {
      // Check voyage info
      await expect(page.locator('[data-testid="voyage-summary"]')).toContainText('出发日期')
      await expect(page.locator('[data-testid="voyage-summary"]')).toContainText('航线')
      
      // Check cabin info
      await expect(page.locator('[data-testid="cabin-summary"]')).toContainText('舱房类型')
      
      // Check price breakdown
      await expect(page.locator('[data-testid="price-summary"]')).toContainText('订单总额')
      await expect(page.locator('[data-testid="total-amount"]')).toBeVisible()
    })

    test('should require terms agreement before submit', async ({ page }) => {
      // Try to submit without checking terms
      const submitButton = page.locator('button:has-text("确认并支付")')
      await expect(submitButton).toBeDisabled()
      
      // Check terms checkbox
      await page.check('input[name="agreeTerms"]')
      await expect(submitButton).toBeEnabled()
    })

    test('should create order and redirect to payment page', async ({ page }) => {
      // Mock API response
      await page.route('**/api/orders', async route => {
        await route.fulfill({
          status: 201,
          body: JSON.stringify({
            data: {
              id: 'order-123',
              order_number: 'ORD202401010001',
              total_amount: 2950
            }
          })
        })
      })
      
      // Check terms and submit
      await page.check('input[name="agreeTerms"]')
      await page.click('button:has-text("确认并支付")')
      
      // Should redirect to payment page
      await page.waitForURL('/payment/order-123')
      await expect(page.locator('h1')).toContainText('订单支付')
    })

    test('should show error if order creation fails', async ({ page }) => {
      // Mock API error
      await page.route('**/api/orders', async route => {
        await route.fulfill({
          status: 400,
          body: JSON.stringify({
            message: '舱房库存不足'
          })
        })
      })
      
      // Check terms and submit
      await page.check('input[name="agreeTerms"]')
      await page.click('button:has-text("确认并支付")')
      
      // Should show error message
      await expect(page.locator('text=舱房库存不足')).toBeVisible()
    })
  })

  test.describe('Payment Page', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/payment/order-123')
      await page.waitForSelector('[data-testid="payment-page"]')
    })

    test('should display payment page with order details', async ({ page }) => {
      // Check order info
      await expect(page.locator('[data-testid="order-summary"]')).toBeVisible()
      await expect(page.locator('[data-testid="order-number"]')).toContainText('ORD')
      await expect(page.locator('[data-testid="pay-amount"]')).toBeVisible()
      
      // Check payment methods
      await expect(page.locator('[data-testid="payment-methods"]')).toBeVisible()
      await expect(page.locator('text=微信支付')).toBeVisible()
    })

    test('should select payment method', async ({ page }) => {
      // Click on WeChat Pay
      await page.click('[data-testid="payment-method-wechat"]')
      
      // Check it's selected
      await expect(page.locator('[data-testid="payment-method-wechat"].selected')).toBeVisible()
    })

    test('should show QR code for WeChat payment', async ({ page }) => {
      // Mock payment creation
      await page.route('**/api/payments', async route => {
        await route.fulfill({
          status: 200,
          body: JSON.stringify({
            data: {
              id: 'payment-123',
              pay_url: 'https://api.example.com/qr-code.png',
              payment_no: 'PAY20240101123456'
            }
          })
        })
      })
      
      // Select payment method and click pay
      await page.click('[data-testid="payment-method-wechat"]')
      await page.click('button:has-text("立即支付")')
      
      // Should show QR code
      await page.waitForSelector('[data-testid="qr-code"]')
      await expect(page.locator('img[alt="支付二维码"]')).toBeVisible()
    })

    test('should poll for payment status', async ({ page }) => {
      // Mock payment status check
      let pollCount = 0
      await page.route('**/api/payments/**', async route => {
        pollCount++
        await route.fulfill({
          status: 200,
          body: JSON.stringify({
            data: {
              id: 'payment-123',
              status: pollCount >= 2 ? 'success' : 'pending'
            }
          })
        })
      })
      
      await page.goto('/payment/order-123')
      
      // Wait for polling to complete
      await page.waitForTimeout(6000)
      
      // Should redirect to success page
      await page.waitForURL('/payment/result?order_id=order-123&status=success')
    })

    test('should allow cancelling order', async ({ page }) => {
      // Mock cancel API
      await page.route('**/api/orders/order-123/cancel', async route => {
        await route.fulfill({
          status: 200,
          body: JSON.stringify({ data: { success: true } })
        })
      })
      
      // Click cancel
      await page.click('text=取消订单')
      
      // Confirm in dialog
      page.once('dialog', async dialog => {
        await dialog.accept()
      })
      
      // Should redirect to orders page
      await page.waitForURL('/orders')
    })

    test('should show countdown timer', async ({ page }) => {
      // Check countdown is visible
      await expect(page.locator('[data-testid="countdown"]')).toBeVisible()
      
      // Verify it shows time in format MM:SS
      const countdownText = await page.locator('[data-testid="countdown"]').textContent()
      expect(countdownText).toMatch(/\d{2}:\d{2}/)
    })
  })

  test.describe('Payment Result Page', () => {
    test('should display success page for successful payment', async ({ page }) => {
      await page.goto('/payment/result?order_id=order-123&status=success')
      
      // Check success message
      await expect(page.locator('h1')).toContainText('支付成功')
      await expect(page.locator('text=您的订单已确认')).toBeVisible()
      
      // Check order details
      await expect(page.locator('[data-testid="order-info"]')).toBeVisible()
      
      // Check action buttons
      await expect(page.locator('a:has-text("查看订单详情")')).toBeVisible()
      await expect(page.locator('a:has-text("返回首页")')).toBeVisible()
    })

    test('should display failed page for failed payment', async ({ page }) => {
      await page.goto('/payment/result?order_id=order-123&status=failed')
      
      // Check failure message
      await expect(page.locator('h1')).toContainText('支付失败')
      await expect(page.locator('text=请重新尝试支付')).toBeVisible()
      
      // Check action buttons
      await expect(page.locator('a:has-text("重新支付")')).toBeVisible()
    })

    test('should navigate to order detail from success page', async ({ page }) => {
      await page.goto('/payment/result?order_id=order-123&status=success')
      
      // Click view order
      await page.click('a:has-text("查看订单详情")')
      
      // Should navigate to order detail
      await page.waitForURL('/orders/order-123')
    })
  })

  test.describe('Full Booking Flow', () => {
    test('should complete end-to-end booking flow', async ({ page }) => {
      // Step 1: Select voyage
      await page.goto('/cruises/test-cruise-id')
      await page.click('text=立即预订')
      await page.waitForURL('/booking**')
      
      await page.click('[data-testid="voyage-card"]:first-child')
      await page.click('button:has-text("下一步")')
      
      // Step 2: Select cabin
      await page.waitForSelector('[data-testid="cabin-selection"]')
      await page.click('[data-testid="cabin-card"]:first-child button:has-text("选择此舱房")')
      await page.click('button:has-text("下一步")')
      
      // Step 3: Fill passenger info
      await page.waitForSelector('[data-testid="passenger-form"]')
      await page.fill('input[name="contactName"]', '张三')
      await page.fill('input[name="contactPhone"]', '13800138000')
      await page.fill('[data-testid="passenger-card"]:first-child input[name="surname"]', 'ZHANG')
      await page.fill('[data-testid="passenger-card"]:first-child input[name="givenName"]', 'SAN')
      await page.fill('[data-testid="passenger-card"]:first-child input[name="name"]', '张三')
      await page.selectOption('[data-testid="passenger-card"]:first-child select[name="gender"]', 'male')
      await page.fill('[data-testid="passenger-card"]:first-child input[name="birthDate"]', '1990-01-01')
      await page.click('button:has-text("下一步")')
      
      // Step 4: Confirm order
      await page.waitForSelector('[data-testid="order-confirmation"]')
      await page.check('input[name="agreeTerms"]')
      
      // Mock order creation
      await page.route('**/api/orders', async route => {
        await route.fulfill({
          status: 201,
          body: JSON.stringify({
            data: { id: 'order-123', order_number: 'ORD202401010001', total_amount: 2950 }
          })
        })
      })
      
      await page.click('button:has-text("确认并支付")')
      
      // Payment page
      await page.waitForURL('/payment/order-123')
      await page.click('[data-testid="payment-method-wechat"]')
      
      // Mock payment
      await page.route('**/api/payments', async route => {
        await route.fulfill({
          status: 200,
          body: JSON.stringify({
            data: { id: 'payment-123', pay_url: 'https://example.com/qr.png' }
          })
        })
      })
      
      await page.click('button:has-text("立即支付")')
      
      // Should show QR code
      await page.waitForSelector('[data-testid="qr-code"]')
    })

    test('should handle browser back button correctly', async ({ page }) => {
      // Navigate through steps
      await page.goto('/booking?cruiseId=test-cruise-id')
      await page.click('[data-testid="voyage-card"]:first-child')
      await page.click('button:has-text("下一步")')
      
      await page.waitForSelector('[data-testid="cabin-selection"]')
      await page.click('[data-testid="cabin-card"]:first-child button:has-text("选择此舱房")')
      await page.click('button:has-text("下一步")')
      
      await page.waitForSelector('[data-testid="passenger-form"]')
      
      // Click browser back
      await page.goBack()
      
      // Should be on cabin selection
      await expect(page.locator('[data-testid="cabin-selection"]')).toBeVisible()
      
      // Selected cabin should be preserved
      await expect(page.locator('[data-testid="cabin-card"]:first-child button')).toContainText('已选择')
    })

    test('should preserve form data on navigation', async ({ page }) => {
      await page.goto('/booking?cruiseId=test-cruise-id&step=3')
      
      // Fill form
      await page.fill('input[name="contactName"]', '张三')
      await page.fill('input[name="contactPhone"]', '13800138000')
      
      // Navigate back and forth
      await page.click('button:has-text("上一步")')
      await page.click('button:has-text("下一步")')
      
      // Form data should be preserved
      await expect(page.locator('input[name="contactName"]')).toHaveValue('张三')
      await expect(page.locator('input[name="contactPhone"]')).toHaveValue('13800138000')
    })
  })

  test.describe('Error Handling', () => {
    test('should show error when voyage loading fails', async ({ page }) => {
      // Mock API error
      await page.route('**/api/voyages**', async route => {
        await route.fulfill({ status: 500, body: 'Internal Server Error' })
      })
      
      await page.goto('/booking?cruiseId=test-cruise-id')
      
      // Should show error message
      await expect(page.locator('text=加载失败')).toBeVisible()
    })

    test('should show error when cabin inventory is insufficient', async ({ page }) => {
      await page.goto('/booking?cruiseId=test-cruise-id&voyageId=test-voyage-id')
      
      // Mock cabin API with 0 available
      await page.route('**/api/cabins**', async route => {
        await route.fulfill({
          status: 200,
          body: JSON.stringify({
            data: [{
              id: 'cabin-1',
              cabin_type: { name_cn: '内舱房' },
              inventory: { available_cabins: 0 }
            }]
          })
        })
      })
      
      await page.waitForSelector('[data-testid="cabin-list"]')
      
      // Select button should be disabled
      const selectButton = page.locator('[data-testid="cabin-card"]:first-child button')
      await expect(selectButton).toBeDisabled()
    })

    test('should handle network timeout during order creation', async ({ page }) => {
      await page.goto('/booking?cruiseId=test-cruise-id&step=4')
      await page.check('input[name="agreeTerms"]')
      
      // Mock timeout
      await page.route('**/api/orders', async route => {
        await new Promise(resolve => setTimeout(resolve, 30000))
        await route.fulfill({ status: 408 })
      })
      
      await page.click('button:has-text("确认并支付")')
      
      // Should show timeout error
      await expect(page.locator('text=请求超时')).toBeVisible()
    })
  })

  test.describe('Responsive Design', () => {
    test('should display correctly on mobile viewport', async ({ page }) => {
      // Set mobile viewport
      await page.setViewportSize({ width: 375, height: 667 })
      
      await page.goto('/booking?cruiseId=test-cruise-id')
      await page.waitForSelector('[data-testid="voyage-list"]')
      
      // Check mobile layout
      await expect(page.locator('[data-testid="booking-wizard"]')).toBeVisible()
      
      // Cards should be full width on mobile
      const voyageCard = page.locator('[data-testid="voyage-card"]').first()
      const box = await voyageCard.boundingBox()
      expect(box.width).toBeLessThanOrEqual(375)
    })

    test('should display correctly on tablet viewport', async ({ page }) => {
      await page.setViewportSize({ width: 768, height: 1024 })
      
      await page.goto('/booking?cruiseId=test-cruise-id')
      await page.waitForSelector('[data-testid="voyage-list"]')
      
      // Layout should adapt to tablet
      await expect(page.locator('[data-testid="voyage-list"]')).toBeVisible()
    })
  })
})
