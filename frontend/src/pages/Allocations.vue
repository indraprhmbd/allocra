<script setup lang="ts">
import { ref, onMounted } from "vue";
import {
  IconPlus,
  IconAlertCircle,
  IconChevronDown,
  IconChevronUp,
} from "@tabler/icons-vue";
import DataTable, { type Column } from "../components/DataTable.vue";
import StatusBadge from "../components/StatusBadge.vue";
import ModalBase from "../components/ModalBase.vue";
import api from "../services/api";
import { useToast } from "../composables/useToast";

const { toast } = useToast();

const showCreateModal = ref(false);
const expandedRows = ref<number[]>([]);
const loading = ref(false);

const columns: Column[] = [
  { key: "id", label: "ID", class: "w-24 font-mono text-xs" },
  { key: "room_id", label: "Resource ID", class: "font-mono" },
  { key: "timeRange", label: "Time Scope" },
  { key: "status", label: "Execution_Status", class: "text-right" },
  { key: "expand", label: "", class: "w-10" },
];

const requests = ref<any[]>([]);

const forceAllocate = async (booking: any) => {
  // Assuming 'Booking' type is 'any' for now
  try {
    await api.patch(`/bookings/${booking.id}/force`);
    toast(`Node ${booking.room_id} preempted & reallocated`, "success");
    fetchRequests();
  } catch (err) {
    toast("Preemption failed", "error");
  }
};

const handleAction = async (id: number, action: "approve" | "reject") => {
  try {
    await api.patch(`/bookings/${id}/${action}`);
    toast(`Allocation ${action}d successfully`, "success");
    fetchRequests();
  } catch (err) {
    toast(`Action failed: ${action}`, "error");
  }
};

const fetchRequests = async () => {
  loading.value = true;
  try {
    const res = await api.get("/bookings/all");
    requests.value = res.data.map((b: any) => ({
      ...b,
      displayId: `AL-${b.id}`,
      timeRange: `${new Date(b.start_time).toLocaleString()} - ${new Date(b.end_time).toLocaleTimeString()}`,
    }));
  } catch (err) {
    console.error("Failed to fetch allocations");
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchRequests();
});

const toggleExpand = (id: number) => {
  const index = expandedRows.value.indexOf(id);
  if (index > -1) expandedRows.value.splice(index, 1);
  else expandedRows.value.push(id);
};

const isExpanded = (id: number) => expandedRows.value.includes(id);
</script>

<template>
  <div class="space-y-8">
    <div class="flex items-end justify-between">
      <div>
        <h2 class="text-2xl font-bold text-primary tracking-tight">
          Allocation Requests
        </h2>
        <p class="text-muted text-sm mt-1">
          Review pending resource reservations and resolve engine-detected
          conflicts.
        </p>
      </div>
      <button
        @click="showCreateModal = true"
        class="bg-accent hover:bg-accent-hover text-white text-xs font-bold uppercase tracking-widest px-4 py-2 rounded-sm transition-colors flex items-center gap-2"
      >
        <IconPlus :size="16" />
        NEW_ALLOCATION
      </button>
    </div>

    <!-- Table -->
    <div class="bg-surface border border-border rounded-sm overflow-hidden">
      <DataTable :columns="columns" :items="requests">
        <template #id="{ item }">
          <span class="text-accent font-mono">{{ item.displayId }}</span>
        </template>

        <template #status="{ value }">
          <StatusBadge
            :status="value"
            :variant="
              value === 'allocated'
                ? 'info'
                : value === 'pending'
                  ? 'default'
                  : 'error'
            "
          />
        </template>

        <template #expand="{ item }">
          <button
            v-if="item.conflict"
            @click="toggleExpand(item.id)"
            class="text-muted hover:text-primary p-1"
          >
            <IconChevronUp v-if="isExpanded(item.id)" :size="16" />
            <IconChevronDown v-else :size="16" />
          </button>
        </template>
      </DataTable>

      <!-- Custom Expandable Section Logic outside of DataTable shared logic if needed, 
           or via a custom template in DataTable that spans full width. 
           For this high-precision UI, we'll keep it inline in the requests view. -->
      <template v-for="item in requests" :key="'conflict-' + item.id">
        <div
          v-if="isExpanded(item.id)"
          class="bg-red-500/5 border-b border-border p-6 flex gap-6 items-start animate-fade"
        >
          <IconAlertCircle :size="24" class="text-red-500 shrink-0" />
          <div class="space-y-3">
            <h4
              class="text-xs font-bold text-red-500 uppercase tracking-widest"
            >
              Conflict Detected: {{ item.conflict?.type }}
            </h4>
            <p class="text-sm text-muted max-w-2xl">
              This request overlaps with active allocation
              <span class="text-primary font-mono">{{
                item.conflict?.overlapping_id
              }}</span
              >.
              {{ item.conflict?.details }}
            </p>
            <div class="flex gap-4 pt-2">
              <button
                @click="forceAllocate(item)"
                class="flex-1 bg-accent hover:bg-accent-hover text-white text-[10px] font-bold uppercase py-2.5 rounded-sm transition-colors"
              >
                FORCE_ALLOCATE
              </button>
              <button
                @click="handleAction(item.id, 'reject')"
                class="flex-1 border border-border hover:bg-white/5 text-muted hover:text-primary text-[10px] font-bold uppercase py-2.5 rounded-sm transition-colors"
              >
                CONFIRM_REJECTION
              </button>
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- Create Modal -->
    <ModalBase
      :show="showCreateModal"
      title="New Allocation Request"
      @close="showCreateModal = false"
    >
      <AllocationForm
        @cancel="showCreateModal = false"
        @submit="showCreateModal = false"
      />
    </ModalBase>
  </div>
</template>

<style scoped>
.animate-fade {
  animation: fadeIn 0.15s ease-out;
}
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
