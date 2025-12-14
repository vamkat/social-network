import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { apiRequest } from '@/lib/api'

export const useStore = create(
  persist(
    (set) => ({
      // State
      user: null,
      loading: false,

      // Fetch user profile by ID
      loadUserProfile: async (userId) => {
        set({ loading: true })
        try {
          console.log("calling backend for user data")
          const userData = await apiRequest(`/profile/${userId}`, {
            method: "POST",
          })
          set({ user: { id: userData.user_id, avatar: userData.avatar }, loading: false })
          return { success: true }
        } catch (error) {
          console.error('Failed to load user profile:', error)
          set({ user: null, loading: false })
          return { success: false, error: error.message }
        }
      },

      // Clear user (on logout)
      clearUser: () => {
        set({ user: null })
      },
    }),
    {
      name: 'user',
      partialize: (state) => ({
        user: state.user
      }),
    }
  )
)