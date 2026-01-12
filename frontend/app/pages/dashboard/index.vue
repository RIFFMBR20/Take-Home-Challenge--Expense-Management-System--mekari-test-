<template>
  <div class="min-h-screen bg-gray-100 text-gray-900">
    <div class="bg-white border-b border-gray-200 p-6 shadow-sm">
      <div class="max-w-6xl mx-auto flex justify-between items-center">
        <div>
          <h2 class="text-2xl font-bold text-gray-800">
            Expense Management
          </h2>
          <span class="bg-blue-100 text-blue-700 px-3 py-1 rounded-full text-xs font-bold uppercase mt-1 inline-block">
            Login sebagai: {{ userRole || 'User' }}
          </span>
        </div>
        <div class="flex gap-3">
          <button
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md transition font-medium shadow-sm"
            @click="showModal = true"
          >
            + Tambah Pengeluaran
          </button>
          <button
            class="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-md transition font-medium"
            @click="logout"
          >
            Keluar
          </button>
        </div>
      </div>
    </div>

    <div class="p-8">
      <div class="max-w-6xl mx-auto bg-white rounded-lg shadow-md overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-gray-50 border-b border-gray-200 text-gray-600 uppercase text-sm">
                <th class="p-4 font-semibold">
                  ID
                </th>
                <th class="p-4 font-semibold">
                  Deskripsi
                </th>
                <th class="p-4 font-semibold">
                  Nominal
                </th>
                <th class="p-4 font-semibold">
                  Bukti
                </th>
                <th class="p-4 font-semibold text-center">
                  Aksi
                </th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="item in expenseList"
                :key="item.id"
                class="border-b border-gray-100 hover:bg-gray-50 transition"
              >
                <td class="p-4 text-gray-500 text-sm">
                  #{{ item.id }}
                </td>
                <td class="p-4 font-medium">
                  {{ item.description }}
                </td>
                <td class="p-4 font-mono text-blue-600 font-bold">
                  {{ formatRupiah(item.amount) }}
                </td>
                <td class="p-4">
                  <a
                    v-if="item.receipt_url"
                    :href="item.receipt_url"
                    target="_blank"
                    class="text-blue-500 underline text-xs"
                  >Lihat Bukti</a>
                  <span
                    v-else
                    class="text-gray-400 text-xs italic"
                  >Tanpa Bukti</span>
                </td>
                <td class="p-4 text-center">
                  <div
                    v-if="userRole === 'manager' && item.status === 'pending'"
                    class="flex justify-center gap-2"
                  >
                    <button
                      class="bg-green-500 text-white px-2 py-1 rounded text-xs hover:bg-green-600"
                      @click="updateStatus(item.id, 'approved')"
                    >
                      Approve
                    </button>
                    <button
                      class="bg-red-500 text-white px-2 py-1 rounded text-xs hover:bg-red-600"
                      @click="updateStatus(item.id, 'rejected')"
                    >
                      Reject
                    </button>
                  </div>
                  <span
                    v-else
                    class="px-2 py-1 rounded text-[10px] font-bold uppercase"
                    :class="getStatusClass(item.status)"
                  >
                    {{ item.status }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <div
      v-if="showModal"
      class="fixed inset-0 bg-black bg-opacity-30 backdrop-blur-sm flex items-center justify-center z-50 p-4"
    >
      <div class="bg-white rounded-lg shadow-2xl w-full max-w-md p-6">
        <h3 class="text-xl font-bold mb-4 text-gray-800 border-b pb-2">
          Tambah Pengeluaran
        </h3>
        <form @submit.prevent="submitExpense">
          <div class="mb-4">
            <label class="block text-sm font-medium mb-1">Deskripsi</label>
            <input
              v-model="form.description"
              type="text"
              class="w-full border rounded-md p-2 outline-none focus:ring-2 focus:ring-blue-500"
              required
            >
          </div>
          <div class="mb-4">
            <label class="block text-sm font-medium mb-1">Nominal (IDR)</label>
            <input
              v-model.number="form.amount"
              type="number"
              class="w-full border rounded-md p-2 outline-none focus:ring-2 focus:ring-blue-500"
              required
            >
          </div>
          <div class="mb-6">
            <label class="block text-sm font-medium mb-1">Bukti Pembayaran <span class="text-gray-400 font-normal">(Opsional)</span></label>
            <input
              type="file"
              class="w-full text-xs text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-xs file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
              @change="handleFileUpload"
            >
          </div>
          <div class="flex justify-end gap-2">
            <button
              type="button"
              class="bg-gray-100 hover:bg-gray-200 px-4 py-2 rounded-md"
              @click="showModal = false"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="submitting"
              class="bg-blue-600 text-white px-4 py-2 rounded-md disabled:opacity-50"
            >
              {{ submitting ? 'Menyimpan...' : 'Simpan' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
const expenseList = ref([])
const tokenCookie = useCookie('auth_token')
const userRole = useCookie('user_role')
const showModal = ref(false)
const submitting = ref(false)

const form = ref({
  description: '',
  amount: 0,
  receipt: null // Untuk menyimpan file bukti
})

const handleFileUpload = (event) => {
  const file = event.target.files[0]
  form.value.receipt = file
}

const submitExpense = async () => {
  // Validasi minimal nominal di frontend agar tidak kena error 422 dari Go
  if (form.value.amount < 10000) {
    alert('Minimal pengeluaran adalah Rp 10.000')
    return
  }

  submitting.value = true
  try {
    // KIRIM SEBAGAI JSON (Bukan FormData)
    const payload = {
      description: form.value.description,
      amount: Number(form.value.amount), // Wajib number agar tidak error Decode di Go
      user_id: 1 // Sesuaikan dengan logika ID user kamu jika perlu
    }

    await $fetch('http://localhost:8080/api/expenses', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${tokenCookie.value}`,
        'Content-Type': 'application/json'
      },
      body: payload
    })

    alert('Pengeluaran berhasil disimpan!')
    showModal.value = false
    form.value = { description: '', amount: 0, receipt: null } // Reset
    await fetchExpenses() // Refresh tabel
  } catch (err) {
    const errorMsg = err.data?.message || err.message
    alert('Gagal simpan: ' + errorMsg)
  } finally {
    submitting.value = false
  }
}

const fetchExpenses = async () => {
  try {
    const res = await $fetch('http://localhost:8080/api/expenses', {
      headers: { Authorization: `Bearer ${tokenCookie.value}` }
    })
    if (res && res.data) expenseList.value = res.data
  } catch (err) { console.error(err) }
}

const getStatusClass = s => s === 'pending' ? 'bg-orange-100 text-orange-600' : s === 'rejected' ? 'bg-red-100 text-red-600' : 'bg-green-100 text-green-600'
const formatRupiah = v => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v)
const logout = () => { tokenCookie.value = null; userRole.value = null; window.location.href = '/' }
const updateStatus = async (id, s) => { /* logic update status */ }

onMounted(fetchExpenses)
</script>
