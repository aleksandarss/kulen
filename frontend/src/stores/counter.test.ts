import { setActivePinia, createPinia } from 'pinia'
import { useCounterStore } from './counter'
import { describe, it, expect, beforeEach } from 'vitest'

describe('counter store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('increments', () => {
    const store = useCounterStore()
    store.increment()
    expect(store.count).toBe(1)
    expect(store.doubleCount).toBe(2)
  })
})
