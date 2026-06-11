import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { pb } from '@/boot/pocketbase';

export const useAuthStore = defineStore('auth', () => {
  const user = ref(pb.authStore.record);
  // Configure as a reactive ref to guarantee Vue-Router tracks it correctly
  const isLoggedIn = ref(pb.authStore.isValid);
  const isAdmin = computed(() => !!user.value?.is_admin);

  // Sync state changes from PocketBase's internal authStore
  pb.authStore.onChange((_token, model) => {
    user.value = model;
    isLoggedIn.value = pb.authStore.isValid;
  });

  async function loginEmail(email: string, pass: string) {
    await pb.collection('users').authWithPassword(email, pass);
    user.value = pb.authStore.record;
    isLoggedIn.value = pb.authStore.isValid; // Synchronously force updates
  }

  async function register(email: string, pass: string, name: string) {
    await pb.collection('users').create({
      email,
      password: pass,
      passwordConfirm: pass,
      username: name,
    });
    await loginEmail(email, pass);
  }

  // Google OAuth2 handler
  async function loginGoogle() {
    await pb.collection('users').authWithOAuth2({ provider: 'google' });
    user.value = pb.authStore.record;
    isLoggedIn.value = pb.authStore.isValid; // Synchronously force updates
  }

  function logout() {
    pb.authStore.clear();
    user.value = null;
    isLoggedIn.value = false; // Synchronously clear state
  }

  return {
    user,
    isLoggedIn,
    isAdmin,
    loginEmail,
    register,
    loginGoogle,
    logout,
  };
});
