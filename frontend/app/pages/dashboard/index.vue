<template>
  <div class="min-h-screen bg-gray-100 text-gray-900">
    <div class="bg-white border-b border-gray-200 p-6 shadow-sm">
      <div class="max-w-6xl mx-auto flex justify-between items-center">
        <div>
          <h2 class="text-2xl font-bold text-gray-800">
            Expense Management
          </h2>
          <div class="flex items-center gap-2 mt-1">
            <span class="bg-blue-100 text-blue-700 px-3 py-1 rounded-full text-xs font-bold uppercase">
              Role: {{ userRole }}
            </span>
          </div>
        </div>
        <div class="flex gap-3">
          <button
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md transition font-medium shadow-sm"
            @click="showModal = true"
          >
            + Tambah
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
      <div class="max-w-6xl mx-auto">
        <div class="mb-4 flex justify-end">
          <select
            v-model="filterStatus"
            class="border rounded-md p-2 text-sm outline-none focus:ring-2 focus:ring-blue-500"
            @change="fetchExpenses(1)"
          >
            <option value="">
              Semua Status
            </option>
            <option value="pending">
              Pending
            </option>
            <option value="auto_approved">
              Auto Approved
            </option>
            <option value="completed">
              Completed
            </option>
            <option value="rejected">
              Rejected
            </option>
          </select>
        </div>

        <div class="bg-white rounded-lg shadow-md overflow-hidden">
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
                    class="text-blue-500 underline text-xs font-bold"
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
                      class="bg-green-500 text-white px-2 py-1 rounded text-xs hover:bg-green-600 font-bold"
                      @click="updateStatus(item.id, 'approve')"
                    >
                      Approve
                    </button>
                    <button
                      class="bg-red-500 text-white px-2 py-1 rounded text-xs hover:bg-red-600 font-bold"
                      @click="updateStatus(item.id, 'reject')"
                    >
                      Reject
                    </button>
                  </div>
                  <span
                    v-else
                    :class="getStatusClass(item.status)"
                    class="px-2 py-1 rounded text-[10px] font-bold uppercase inline-block"
                  >
                    {{ item.status.replace('_', ' ') }}
                  </span>
                </td>
              </tr>
              <tr v-if="expenseList.length === 0">
                <td
                  colspan="5"
                  class="p-10 text-center text-gray-400 italic"
                >
                  Data tidak ditemukan.
                </td>
              </tr>
            </tbody>
          </table>

          <div class="bg-gray-50 px-6 py-4 border-t border-gray-200 flex items-center justify-between">
            <p class="text-sm text-gray-600">
              Total: <span class="font-bold">{{ metaData.total }}</span> pengeluaran
            </p>
            <div class="flex gap-2">
              <button
                :disabled="currentPage === 1"
                class="px-3 py-1 border rounded bg-white disabled:opacity-50 text-sm"
                @click="changePage(currentPage - 1)"
              >
                Prev
              </button>
              <span class="flex items-center px-3 text-sm font-bold">Hal {{ currentPage }}</span>
              <button
                :disabled="currentPage >= Math.ceil(metaData.total / metaData.limit)"
                class="px-3 py-1 border rounded bg-white disabled:opacity-50 text-sm"
                @click="changePage(currentPage + 1)"
              >
                Next
              </button>
            </div>
          </div>
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
              class="w-full border rounded-md p-2 outline-none focus:ring-2 focus:ring-blue-500 text-sm"
              required
            >
          </div>
          <div class="mb-4">
            <label class="block text-sm font-medium mb-1">Nominal (Maks 50jt)</label>
            <input
              v-model.number="form.amount"
              type="number"
              :class="{ 'border-red-500': form.amount > 50000000 }"
              class="w-full border rounded-md p-2 outline-none focus:ring-2 focus:ring-blue-500 text-sm"
              required
            >
            <p
              v-if="form.amount > 50000000"
              class="text-red-500 text-xs mt-1"
            >
              Melebihi batas maksimal Rp 50.000.000
            </p>
          </div>
          <div class="mb-6">
            <label class="block text-sm font-medium mb-1 text-xs">Bukti (Opsional)</label>
            <input
              type="file"
              class="w-full text-xs"
              @change="handleFileUpload"
            >
          </div>
          <div class="flex justify-end gap-2">
            <button
              type="button"
              class="bg-gray-100 px-4 py-2 rounded-md text-sm"
              @click="showModal = false"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="submitting || form.amount > 50000000"
              class="bg-blue-600 text-white px-4 py-2 rounded-md text-sm disabled:opacity-50"
            >
              {{ submitting ? 'Proses...' : 'Simpan' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
const expenseList = ref([])
const metaData = ref({ total: 0, page: 1, limit: 10 })
const currentPage = ref(1)
const filterStatus = ref('')
const showModal = ref(false)
const submitting = ref(false)

const tokenCookie = useCookie('auth_token')
const userId = useCookie('user_id')
const userRole = useCookie('user_role')

const form = ref({ description: '', amount: 0, receipt: null, receipt_url: null })

const fetchExpenses = async (page = 1) => {
  try {
    const url = `http://localhost:8080/api/expenses?page=${page}&limit=10&status=${filterStatus.value}`

    const res = await $fetch(url, {
      headers: {
        'Authorization': `Bearer ${tokenCookie.value}`,
        'X-User-ID': userId.value || 1,
        'X-User-Role': userRole.value || 'employee'
      }
    })

    if (res) {
      expenseList.value = res.data || []
      metaData.value = res.meta || { total: 0, page: 1, limit: 10 }
      currentPage.value = res.meta?.page || page
    }
  } catch (err) {
    console.error('Fetch Error:', err)
  }
}

const changePage = (newPage) => {
  const maxPage = Math.ceil(metaData.value.total / metaData.value.limit)
  if (newPage > 0 && newPage <= maxPage) fetchExpenses(newPage)
}

const handleFileUpload = (event) => {
  const file = event.target.files[0]
  if (file) {
    form.value.receipt_url = URL.createObjectURL(file)
    form.value.receipt = file
  }
}

const submitExpense = async () => {
  if (form.value.amount < 10000) return alert('Min Rp 10.000')
  submitting.value = true
  try {
    await $fetch('http://localhost:8080/api/expenses', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${tokenCookie.value}`,
        'Content-Type': 'application/json'
      },
      body: {
        description: form.value.description,
        amount: Number(form.value.amount)
      }
    })
    showModal.value = false
    form.value = { description: '', amount: 0, receipt: null, receipt_url: null }
    await fetchExpenses(1)
  } catch (err) {
    alert(err.data?.error || 'Gagal menyimpan')
  } finally {
    submitting.value = false
  }
}

const updateStatus = async (id, action) => {
  if (!confirm(`Konfirmasi ${action}?`)) return
  try {
    await $fetch(`http://localhost:8080/api/expenses/${id}/${action}`, {
      method: 'PUT',
      headers: { Authorization: `Bearer ${tokenCookie.value}` }
    })
    await fetchExpenses(currentPage.value)
  } catch (err) {
    alert('Gagal update')
  }
}

const formatRupiah = v => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v)
const getStatusClass = (s) => {
  if (s === 'pending') return 'bg-orange-100 text-orange-600'
  if (s === 'rejected') return 'bg-red-100 text-red-600'
  return 'bg-green-100 text-green-600'
}
const logout = () => {
  tokenCookie.value = null
  userRole.value = null
  window.location.href = '/'
}

onMounted(() => fetchExpenses(1))
</script>
