import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export const useStore = create(
  persist(
    (set) => ({
      // State
      user: null,
      msgReceiver: null,
      loading: false,

      // Manually set user data
      setUser: (userData) => {
        set({ user: userData })
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

export const useMsgReceiver = create(
  persist(
    (set) => ({
      // State
      msgReceiver: null,
      loading: false,

      // set the user data of the user we want to send the msg
      setMsgReceiver: (receiverData) => {
        set({msgReceiver: receiverData})
      },

      // clear receiver
      clearMsgReceiver: () => {
        set({msgReceiver: null})
      },

    }),
    {
      name: 'msgReceiver',
      partialize: (state) => ({
        msgReceiver: state.msgReceiver
      }),
    }
  )
)