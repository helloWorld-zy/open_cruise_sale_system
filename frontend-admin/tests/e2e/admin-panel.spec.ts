import { test, expect } from '@playwright/test'

test.describe('Admin Panel', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/login')
    await page.fill('input[name="username"]', 'admin')
    await page.fill('input[name="password"]', 'password')
    await page.click('button[type="submit"]')
    await page.waitForURL('/')
  })

  test('should display dashboard after login', async ({ page }) => {
    await expect(page.locator('h2')).toContainText('Dashboard')
    await expect(page.locator('aside')).toBeVisible()
  })

  test('should navigate to cruises page', async ({ page }) => {
    await page.click('text=Cruises')
    await page.waitForURL('/cruises')
    await expect(page.locator('h2')).toContainText('邮轮管理')
  })

  test('should create a new cruise', async ({ page }) => {
    await page.goto('/cruises')
    await page.click('text=新增邮轮')
    await page.waitForURL('/cruises/new')

    // Fill form
    await page.selectOption('select[name="companyId"]', '1')
    await page.fill('input[name="nameCn"]', 'Test Cruise Ship')
    await page.fill('input[name="code"]', 'TEST001')
    await page.fill('input[name="grossTonnage"]', '100000')
    await page.fill('input[name="passengerCapacity"]', '5000')

    // Submit
    await page.click('button[type="submit"]')
    await page.waitForURL('/cruises')

    // Verify
    await expect(page.locator('text=Test Cruise Ship')).toBeVisible()
  })

  test('should edit an existing cruise', async ({ page }) => {
    await page.goto('/cruises')
    
    // Click on first cruise's edit button
    await page.click('[data-testid="edit-cruise"]')
    
    // Modify data
    await page.fill('input[name="nameCn"]', 'Updated Cruise Name')
    
    // Save
    await page.click('button[type="submit"]')
    await page.waitForURL('/cruises')
    
    await expect(page.locator('text=Updated Cruise Name')).toBeVisible()
  })

  test('should delete a cruise', async ({ page }) => {
    await page.goto('/cruises')
    
    // Click delete on first cruise
    await page.click('[data-testid="delete-cruise"]')
    
    // Confirm deletion
    await page.click('text=删除')
    
    // Verify success
    await expect(page.locator('text=删除成功')).toBeVisible()
  })

  test('should upload images', async ({ page }) => {
    await page.goto('/cruises/new')
    
    // Trigger file upload
    const fileInput = page.locator('input[type="file"]')
    await fileInput.setInputFiles('tests/fixtures/test-image.jpg')
    
    // Verify preview appears
    await expect(page.locator('img[alt="Preview"]')).toBeVisible()
  })

  test('should manage cabin types', async ({ page }) => {
    await page.goto('/cabin-types')
    
    // Select a cruise
    await page.selectOption('select', { label: 'Test Cruise' })
    
    // Add cabin type
    await page.click('text=新增舱房类型')
    await page.fill('input[name="name"]', 'Balcony Suite')
    await page.fill('input[name="code"]', 'BS001')
    await page.click('button[type="submit"]')
    
    await expect(page.locator('text=Balcony Suite')).toBeVisible()
  })

  test('should manage facilities', async ({ page }) => {
    await page.goto('/facilities')
    
    // Select a cruise
    await page.selectOption('select', { label: 'Test Cruise' })
    
    // Add facility
    await page.click('text=新增设施')
    await page.fill('input[name="name"]', 'Swimming Pool')
    await page.fill('input[name="deckNumber"]', '12')
    await page.click('button[type="submit"]')
    
    await expect(page.locator('text=Swimming Pool')).toBeVisible()
  })
})
