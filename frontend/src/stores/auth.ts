import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { pb } from '@/boot/pocketbase'; // Safe alias path

export const useAuthStore = defineStore('auth', () => {
  const user = ref(pb.authStore.record);
  const isLoggedIn = computed(() => pb.authStore.isValid);

  pb.authStore.onChange((_token, model) => {
    user.value = model;
  });

  async function loginEmail(email: string, pass: string) {
    await pb.collection('users').authWithPassword(email, pass);
    user.value = pb.authStore.record;
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
  }

  function logout() {
    pb.authStore.clear();
    user.value = null;
  }

  return {
    user,
    isLoggedIn,
    loginEmail,
    register,
    loginGoogle,
    logout,
  };
});
