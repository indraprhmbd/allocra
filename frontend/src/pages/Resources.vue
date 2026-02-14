<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import {
  IconPlus,
  IconSearch,
  IconAdjustmentsHorizontal,
} from "@tabler/icons-vue";
import DataTable, { type Column } from "../components/DataTable.vue";
import StatusBadge from "../components/StatusBadge.vue";
import ModalBase from "../components/ModalBase.vue";
import type { Resource } from "../types";
import api from "../services/api";
import { useToast } from "../composables/useToast";

const { toast } = useToast();

const showCreateModal = ref(false);
const showEditModal = ref(false);
const editingId = ref<number | null>(null);
const isEditing = computed(() => editingId.value !== null);
const searchQuery = ref("");
const loading = ref(false);

const columns: Column[] = [
  { key: "name", label: "Identifier", class: "font-mono" },
  { key: "capacity", label: "Capacity", class: "text-right w-32" },
  { key: "type", label: "Access Type", class: "w-40" },
  { key: "status", label: "Status", class: "w-40" },
  { key: "actions", label: "Action", class: "w-20 text-right" },
];

const resources = ref<Resource[]>([]);

const resourceForm = ref({
  name: "",
  capacity: 0,
  type: "shared" as "shared" | "exclusive",
  status: "online",
});

const openEditModal = (resource: Resource) => {
  editingId.value = resource.id;
  resourceForm.value = {
    name: resource.name,
    capacity: resource.capacity,
    type: resource.type as "shared" | "exclusive",
    status: resource.status,
  };
  showEditModal.value = true;
};

const closeModal = () => {
  showCreateModal.value = false;
  showEditModal.value = false;
  editingId.value = null;
  resourceForm.value = {
    name: "",
    capacity: 0,
    type: "shared",
    status: "online",
  };
};

const fetchResources = async () => {
  loading.value = true;
  try {
    const res = await api.get("/rooms");
    resources.value = res.data;
  } catch (err) {
    console.error("Failed to fetch resources");
  } finally {
    loading.value = false;
  }
};

const deleteResource = async (id: number) => {
  if (!confirm("Are you sure you want to decommission this resource?")) return;
  try {
    await api.delete(`/rooms/${id}`);
    toast("Resource successfully decommissioned", "success");
    fetchResources();
  } catch (err) {
    toast("Decommissioning failed", "error");
  }
};

const submitEdit = async () => {
  if (!editingId.value) return;
  try {
    await api.put(`/rooms/${editingId.value}`, resourceForm.value);
    closeModal();
    toast("Resource updated successfully", "success");
    fetchResources();
  } catch (err) {
    toast("Update failed", "error");
  }
};

const submitRegistration = async () => {
  try {
    await api.post("/rooms", resourceForm.value);
    closeModal();
    toast("Resource registered successfully", "success");
    fetchResources();
  } catch (err) {
    toast("Registration failed", "error");
  }
};

onMounted(() => {
  fetchResources();
});
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-primary tracking-tight">
          Resource Management
        </h2>
        <p class="text-muted text-sm mt-1">
          Configure and monitor available system resources and nodes.
        </p>
      </div>
      <button
        @click="showCreateModal = true"
        class="bg-accent hover:bg-accent-hover text-white text-xs font-bold uppercase tracking-widest px-4 py-2 rounded-sm transition-colors flex items-center justify-center gap-2 w-full md:w-auto"
      >
        <IconPlus :size="16" />
        REGISTER_RESOURCE
      </button>
    </div>

    <!-- Toolbar -->
    <div class="flex flex-col md:flex-row items-stretch md:items-center gap-4">
      <div class="relative flex-1 group">
        <IconSearch
          :size="16"
          class="absolute left-3 top-1/2 -translate-y-1/2 text-muted group-focus-within:text-accent transition-colors"
        />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="SEARCH_BY_IDENTIFIER..."
          class="w-full bg-surface border border-border rounded-sm pl-10 pr-4 py-2 text-xs font-mono text-primary placeholder:text-muted focus:outline-none focus:border-accent/50 transition-colors"
        />
      </div>
      <button
        class="border border-border p-2 rounded-sm text-muted hover:text-primary hover:border-muted transition-colors flex justify-center"
      >
        <IconAdjustmentsHorizontal :size="18" />
      </button>
    </div>

    <!-- Table Container -->
    <div class="bg-surface border border-border rounded-sm overflow-hidden">
      <div class="overflow-x-auto">
        <DataTable :columns="columns" :items="resources" :loading="loading">
          <template #type="{ value }">
            <span
              class="text-xs uppercase tracking-wide px-2 py-0.5 border border-border bg-background/50 rounded-full font-mono"
            >
              {{ value }}
            </span>
          </template>

          <template #status="{ value }">
            <StatusBadge
              :status="value"
              :variant="
                value === 'online'
                  ? 'success'
                  : value === 'maintenance'
                    ? 'warning'
                    : 'error'
              "
            />
          </template>

          <template #actions="{ item }">
            <div class="flex items-center justify-end gap-3">
              <button
                @click="openEditModal(item)"
                class="text-[10px] font-mono text-muted hover:text-primary transition-colors"
              >
                EDIT
              </button>
              <button
                @click="deleteResource(item.id)"
                class="text-[10px] font-mono text-red-400/70 hover:text-red-400 transition-colors"
              >
                DELETE
              </button>
            </div>
          </template>
        </DataTable>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <ModalBase
      :show="showCreateModal || showEditModal"
      :title="showEditModal ? 'Update Resource' : 'Register New Resource'"
      @close="closeModal"
    >
      <form
        class="space-y-6"
        @submit.prevent="isEditing ? submitEdit() : submitRegistration()"
      >
        <div class="space-y-2">
          <label
            class="text-[10px] font-mono text-muted uppercase tracking-widest"
            >Resource Identifier</label
          >
          <input
            v-model="resourceForm.name"
            type="text"
            required
            class="w-full bg-background border border-border rounded-sm p-3 text-sm text-primary focus:outline-none focus:border-accent"
            placeholder="e.g. NODE-AX-05"
          />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-2">
            <label
              class="text-[10px] font-mono text-muted uppercase tracking-widest"
              >Capacity Units</label
            >
            <input
              v-model.number="resourceForm.capacity"
              type="number"
              required
              min="1"
              class="w-full bg-background border border-border rounded-sm p-3 text-sm text-primary focus:outline-none focus:border-accent"
              placeholder="0"
            />
          </div>
          <div class="space-y-2">
            <label
              class="text-[10px] font-mono text-muted uppercase tracking-widest"
              >Access Mode</label
            >
            <select
              v-model="resourceForm.type"
              class="w-full bg-background border border-border rounded-sm p-3 text-sm text-primary focus:outline-none focus:border-accent"
            >
              <option value="shared">SHARED</option>
              <option value="exclusive">EXCLUSIVE</option>
            </select>
          </div>
        </div>

        <div v-if="showEditModal" class="space-y-2">
          <label
            class="text-[10px] font-mono text-muted uppercase tracking-widest"
            >Status</label
          >
          <select
            v-model="resourceForm.status"
            class="w-full bg-background border border-border rounded-sm p-3 text-sm text-primary focus:outline-none focus:border-accent"
          >
            <option value="online">ONLINE</option>
            <option value="maintenance">MAINTENANCE</option>
            <option value="offline">OFFLINE</option>
          </select>
        </div>

        <div class="bg-accent/5 border border-accent/20 p-4 rounded-sm">
          <p
            class="text-[10px] text-accent font-medium leading-relaxed uppercase"
          >
            Note: Exclusive access resources can only serve one allocation at a
            time, regardless of requested capacity vs. total units.
          </p>
        </div>
      </form>

      <template #actions>
        <button
          type="button"
          @click="closeModal"
          class="text-xs font-bold font-mono text-muted hover:text-primary px-4 py-2 transition-colors"
        >
          CANCEL
        </button>
        <button
          type="button"
          @click="showEditModal ? submitEdit() : submitRegistration()"
          class="bg-accent hover:bg-accent-hover text-white text-xs font-bold uppercase tracking-widest px-6 py-2 rounded-sm transition-colors"
        >
          {{ showEditModal ? "SAVE_CHANGES" : "SUBMIT_REGISTRATION" }}
        </button>
      </template>
    </ModalBase>
  </div>
</template>
