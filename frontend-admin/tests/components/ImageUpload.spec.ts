import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ImageUpload from '../components/ImageUpload.vue'

describe('ImageUpload', () => {
  it('renders upload area', () => {
    const wrapper = mount(ImageUpload, {
      props: {
        modelValue: []
      }
    })

    expect(wrapper.text()).toContain('点击上传')
    expect(wrapper.text()).toContain('或拖拽文件到此处')
  })

  it('displays preview images when provided', () => {
    const wrapper = mount(ImageUpload, {
      props: {
        modelValue: ['/image1.jpg', '/image2.jpg']
      }
    })

    const images = wrapper.findAll('img')
    expect(images.length).toBe(2)
    expect(images[0].attributes('src')).toBe('/image1.jpg')
    expect(images[1].attributes('src')).toBe('/image2.jpg')
  })

  it('emits remove event when delete button clicked', async () => {
    const wrapper = mount(ImageUpload, {
      props: {
        modelValue: ['/image1.jpg', '/image2.jpg']
      }
    })

    // Find and click delete button on first image
    const deleteButton = wrapper.find('button')
    await deleteButton.trigger('click')

    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([['/image2.jpg']])
  })

  it('accepts multiple prop', () => {
    const wrapper = mount(ImageUpload, {
      props: {
        modelValue: [],
        multiple: true
      }
    })

    const input = wrapper.find('input[type="file"]')
    expect(input.attributes('multiple')).toBeDefined()
  })

  it('accepts max files prop', () => {
    const wrapper = mount(ImageUpload, {
      props: {
        modelValue: [],
        maxFiles: 5
      }
    })

    expect(wrapper.props('maxFiles')).toBe(5)
  })
})
