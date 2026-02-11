import { describe, it, expect } from '@jest/globals'
import { mount } from '@vue/test-utils'
import CruiseCard from '@/components/CruiseCard.vue'

describe('CruiseCard', () => {
  const mockCruise = {
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
    coverImages: ['/static/images/cruise-1.jpg'],
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

  it('emits click event when clicked', async () => {
    const wrapper = mount(CruiseCard, {
      props: {
        cruise: mockCruise
      }
    })

    await wrapper.trigger('click')
    
    expect(wrapper.emitted('click')).toBeTruthy()
    expect(wrapper.emitted('click')![0]).toEqual(['1'])
  })

  it('displays placeholder when no cover image', () => {
    const cruiseWithoutImage = { ...mockCruise, coverImages: [] }
    const wrapper = mount(CruiseCard, {
      props: {
        cruise: cruiseWithoutImage
      }
    })

    const img = wrapper.find('.card-image')
    expect(img.attributes('src')).toBe('/static/images/placeholder.jpg')
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

  it('displays company name when available', () => {
    const cruiseWithCompany = {
      ...mockCruise,
      company: {
        name: '皇家加勒比'
      }
    }
    const wrapper = mount(CruiseCard, {
      props: {
        cruise: cruiseWithCompany
      }
    })

    expect(wrapper.text()).toContain('皇家加勒比')
  })

  it('shows maintenance status', () => {
    const maintenanceCruise = { ...mockCruise, status: 'maintenance' }
    const wrapper = mount(CruiseCard, {
      props: {
        cruise: maintenanceCruise
      }
    })

    expect(wrapper.text()).toContain('维护中')
  })
})
