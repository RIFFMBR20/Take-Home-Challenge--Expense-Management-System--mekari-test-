<template>
  <div class="min-h-screen w-full flex items-center justify-center bg-gray-50 p-4">
    <UCard class="w-full max-w-md shadow-2xl border-t-4 border-blue-600">
      <template #header>
        <div class="text-center py-4">
          <h1 class="text-3xl font-extrabold text-gray-900 tracking-tight">
            Mekari <span class="text-blue-600">Expense</span>
          </h1>
          <p class="mt-2 text-sm text-gray-500">
            Sistem Pengelolaan Pengeluaran
          </p>
        </div>
      </template>

      <form
        class="flex flex-col space-y-5"
        @submit.prevent="handleLogin"
      >
        <UFormGroup
          label="Alamat Email"
          name="email"
        >
          <UInput
            v-model="email"
            type="email"
            placeholder="employee@mekari.com"
            icon="i-heroicons-envelope"
            size="lg"
            block
          />
        </UFormGroup>

        <UFormGroup
          label="Password"
          name="password"
        >
          <UInput
            v-model="password"
            type="password"
            placeholder="********"
            icon="i-heroicons-lock-closed"
            size="lg"
            block
          />
        </UFormGroup>

        <div class="pt-2">
          <UButton
            type="submit"
            block
            size="xl"
            color="blue"
            :loading="loading"
            class="font-bold uppercase tracking-widest"
          >
            Masuk Akun
          </UButton>
        </div>
      </form>

      <template #footer>
        <div class="text-center">
          <p class="text-xs text-gray-400 font-medium">
            © 2026 Mekari Expense Challenge
          </p>
        </div>
      </template>
    </UCard>
  </div>
</template>

<script setup>
const email = ref('')
const password = ref('')
const loading = ref(false)
const toast = useToast()

const handleLogin = async () => {
  loading.value = true
  try {
    const response = await $fetch('http://localhost:8080/api/auth/login', {
      method: 'POST',
      body: { email: email.value, password: password.value }
    })

    // Simpan Cookie
    const token = useCookie('auth_token', { path: '/' })
    const role = useCookie('user_role', { path: '/' })
    const userId = useCookie('user_id', { path: '/' })

    token.value = response.token
    role.value = response.role
    userId.value = response.user_id

    window.location.href = '/dashboard'
  } catch (err) {
    toast.add({ title: 'Gagal Login', description: 'Email atau password salah', color: 'red' })
  } finally {
    loading.value = false
  }
}
</script>
