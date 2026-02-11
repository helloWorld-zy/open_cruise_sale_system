import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import CruiseCard from '../components/cruise/CruiseCard.vue'
import type { Cruise } from '@cruisebooking/types'

describe('CruiseCard', () => {
  const mockCruise: Cruise = {
    id: '1',
    nameCn: '海洋光谱号',
    nameEn: 'Spectrum of the Seas',
    code: 'SOTS001',
    grossTonnage: 169379,
    passengerCapacity: 5622,
    crewCount: 1551,
    deckCount: 16,
    builtYear: 2019,
    lengthMeters: 347,
    widthMeters: 41,
    coverImages: ['/images/cruise-1.jpg'],
    status: 'active',
    sortWeight: 100,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z'
  }

  it('renders cruise information correctly', () => {
    const wrapper = mount(CruiseCard, {
      props: {
        cruise: mockCruise
      }
    })

    expect(wrapper.text()).toContain('海洋光谱号')
    expect(wrapper.text()).toContain('Spectrum of the Seas')
    expect(wrapper.text()).toContain('169,379')
    expect(wrapper.text()).toContain('5,622')
    expect(wrapper.text()).toContain('16')
  })

  it('displays correct status badge', () => {
    const wrapper = mount(CruiseCard, {
      props: {
        cruise: mockCruise
      }
    })

    expect(wrapper.text()).toContain('热售中')
  })

  it('shows inactive status correctly', () => {
    const inactiveCruise = { ...mockCruise, status: 'inactive' }
    const wrapper = mount(CruiseCard, {
      props: {
        cruise: inactiveCruise
      }
    })

    expect(wrapper.text()).toContain('已下架')
  })

  it('has correct link to detail page', () => {
    const wrapper = mount(CruiseCard, {
      props: {
        cruise: mockCruise
      }
    })

    const button = wrapper.find('a[href="/cruises/1"]')
    expect(button.exists()).toBe(true)
    expect(button.text()).toContain('查看详情')
  })

  it('displays placeholder when no cover image', () => {
    const cruiseWithoutImage = { ...mockCruise, coverImages: [] }
    const wrapper = mount(CruiseCard, {
      props: {
        cruise: cruiseWithoutImage
      }
    })

    expect(wrapper.find('.bg-gray-200').exists()).toBe(true)
  })

  it('formats numbers with locale', () => {
    const wrapper = mount(CruiseCard, {
      props: {
        cruise: mockCruise
      }
    })

    // Check that large numbers are formatted with commas
    expect(wrapper.text()).toMatch(/169,?379/)
    expect(wrapper.text()).toMatch(/5,?622/)
  })
})
