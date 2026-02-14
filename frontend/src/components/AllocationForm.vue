<script setup lang="ts">
import { ref, onMounted } from "vue";
import { IconCalendar, IconClock, IconSearch } from "@tabler/icons-vue";
import api from "../services/api";

import { useToast } from "../composables/useToast";

const emit = defineEmits(["submit", "cancel"]);
const { toast } = useToast();

const form = ref({
  room_id: "",
  user_id: 1,
  start_time: "",
  end_time: "",
});

const resources = ref<any[]>([]);

const fetchResources = async () => {
  try {
    const res = await api.get("/rooms");
    resources.value = res.data;
  } catch (err) {
    console.error("Failed to load resources for form");
  }
};

const handleSubmit = async () => {
  if (!form.value.room_id) {
    toast("Please select a resource", "error");
    return;
  }

  try {
    const payload = {
      ...form.value,
      room_id: parseInt(form.value.room_id),
      start_time: new Date(form.value.start_time).toISOString(),
      end_time: new Date(form.value.end_time).toISOString(),
    };
    await api.post("/bookings", payload);
    toast("Allocation request executed", "success");
    emit("submit");
  } catch (err: any) {
    const errorMsg = err.response?.data?.error || "Allocation failed";
    toast(errorMsg, "error");
  }
};

onMounted(() => {
  fetchResources();
});
</script>

<template>
  <form @submit.prevent="handleSubmit" class="space-y-8">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Resource Selection -->
      <div class="md:col-span-2 space-y-2">
        <label
          class="text-[10px] font-mono text-muted uppercase tracking-widest flex items-center gap-2"
        >
          <IconSearch :size="12" /> Target_Resource
        </label>
        <select
          v-model="form.room_id"
          required
          class="w-full bg-background border border-border rounded-sm p-3 text-sm text-primary focus:outline-none focus:border-accent"
        >
          <option value="" disabled>SELECT_RESOURCE_FOR_ALLOCATION...</option>
          <option v-for="r in resources" :key="r.id" :value="r.id">
            {{ r.name }} (MAX_CAPACITY: {{ r.capacity }})
          </option>
        </select>
      </div>

      <!-- Time Range -->
      <div class="space-y-2">
        <label
          class="text-[10px] font-mono text-muted uppercase tracking-widest flex items-center gap-2"
        >
          <IconCalendar :size="12" /> Start_Timestamp
        </label>
        <input
          type="datetime-local"
          v-model="form.start_time"
          required
          class="w-full bg-background border border-border rounded-sm p-3 text-sm text-primary focus:outline-none focus:border-accent font-mono appearance-none"
        />
      </div>

      <div class="space-y-2">
        <label
          class="text-[10px] font-mono text-muted uppercase tracking-widest flex items-center gap-2"
        >
          <IconClock :size="12" /> End_Timestamp
        </label>
        <input
          type="datetime-local"
          v-model="form.end_time"
          required
          class="w-full bg-background border border-border rounded-sm p-3 text-sm text-primary focus:outline-none focus:border-accent font-mono"
        />
      </div>
    </div>

    <!-- Feedback Container -->
    <div
      class="bg-surface p-4 border border-border rounded-sm flex items-start gap-4"
    >
      <div
        class="w-1.5 h-1.5 rounded-full bg-yellow-500 mt-1.5 shadow-[0_0_8px_rgba(234,179,8,0.5)]"
      ></div>
      <p class="text-[10px] font-mono text-muted leading-relaxed">
        ADVISORY: ENGINE_RESOLUTION_MODE IS CURRENTLY SET TO DETERMINISTIC.
        OVERLAPPING ALLOCATIONS WILL BE FLAGGED AS CONFLICTS UNLESS EXPLICITLY
        PERMITTED BY POLICY.
      </p>
    </div>

    <div class="flex justify-end gap-3 pt-4">
      <button
        type="button"
        @click="emit('cancel')"
        class="text-xs font-bold font-mono text-muted hover:text-primary px-6 py-3 transition-colors uppercase"
      >
        DISCARD
      </button>
      <button
        type="submit"
        class="bg-accent hover:bg-accent-hover text-white text-xs font-bold uppercase tracking-widest px-8 py-3 rounded-sm transition-colors shadow-lg shadow-accent/10"
      >
        EXECUTE_ALLOCATION_REQUEST
      </button>
    </div>
  </form>
</template>

<style scoped>
input::-webkit-calendar-picker-indicator {
  filter: invert(1);
}
</style>
