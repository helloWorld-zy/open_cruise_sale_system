import { test, expect } from '@playwright/test'

test.describe('Cruise Browsing', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the cruises page
    await page.goto('/cruises')
    // Wait for the page to load
    await page.waitForSelector('[data-testid="cruise-list"]')
  })

  test('should display cruise list page', async ({ page }) => {
    // Check page title
    await expect(page.locator('h1')).toContainText('邮轮列表')
    
    // Check search filter exists
    await expect(page.locator('input[placeholder="搜索邮轮名称"]')).toBeVisible()
    
    // Check status filter exists
    await expect(page.locator('select')).toBeVisible()
  })

  test('should filter cruises by keyword', async ({ page }) => {
    // Enter search term
    await page.fill('input[placeholder="搜索邮轮名称"]', '光谱号')
    
    // Click search button
    await page.click('button:has-text("搜索")')
    
    // Wait for results
    await page.waitForTimeout(500)
    
    // Check filtered results
    const cruiseCards = page.locator('[data-testid="cruise-card"]')
    await expect(cruiseCards.first()).toContainText('光谱号')
  })

  test('should navigate to cruise detail page', async ({ page }) => {
    // Click on first cruise card
    await page.click('[data-testid="cruise-card"] >> text=查看详情')
    
    // Wait for detail page to load
    await page.waitForURL(/\/cruises\/[\w-]+/)
    
    // Check detail page content
    await expect(page.locator('h1')).toBeVisible()
    await expect(page.locator('[data-testid="cruise-specs"]')).toBeVisible()
  })

  test('should display cruise detail with tabs', async ({ page }) => {
    // Navigate to a specific cruise
    await page.goto('/cruises/test-cruise-id')
    
    // Check all tabs are present
    await expect(page.locator('text=舱房类型')).toBeVisible()
    await expect(page.locator('text=船上设施')).toBeVisible()
    await expect(page.locator('text=邮轮参数')).toBeVisible()
    
    // Check image gallery exists
    await expect(page.locator('[data-testid="image-gallery"]')).toBeVisible()
    
    // Check specs are displayed
    await expect(page.locator('text=总吨位')).toBeVisible()
    await expect(page.locator('text=载客量')).toBeVisible()
  })

  test('should switch between tabs on detail page', async ({ page }) => {
    await page.goto('/cruises/test-cruise-id')
    
    // Click on facilities tab
    await page.click('text=船上设施')
    await expect(page.locator('[data-testid="facility-tabs"]')).toBeVisible()
    
    // Click on specs tab
    await page.click('text=邮轮参数')
    await expect(page.locator('[data-testid="specs-list"]')).toBeVisible()
  })

  test('should handle empty search results', async ({ page }) => {
    // Search for non-existent cruise
    await page.fill('input[placeholder="搜索邮轮名称"]', '不存在的邮轮XYZ123')
    await page.click('button:has-text("搜索")')
    
    // Wait for results
    await page.waitForTimeout(500)
    
    // Check empty state
    await expect(page.locator('text=暂无邮轮')).toBeVisible()
    await expect(page.locator('text=请调整筛选条件或稍后再试')).toBeVisible()
  })

  test('should show loading state', async ({ page }) => {
    // Slow down network to see loading state
    await page.route('**/api/v1/cruises**', async route => {
      await new Promise(r => setTimeout(r, 1000))
      await route.continue()
    })
    
    await page.goto('/cruises')
    
    // Check loading indicator is shown
    await expect(page.locator('text=加载中...')).toBeVisible()
  })

  test('should navigate back from detail page', async ({ page }) => {
    await page.goto('/cruises')
    const initialUrl = page.url()
    
    // Go to detail
    await page.click('[data-testid="cruise-card"] >> text=查看详情')
    await page.waitForURL(/\/cruises\/[\w-]+/)
    
    // Go back
    await page.click('text=返回列表')
    await page.waitForURL(initialUrl)
    
    // Verify we're back on list page
    await expect(page.locator('h1')).toContainText('邮轮列表')
  })
})
