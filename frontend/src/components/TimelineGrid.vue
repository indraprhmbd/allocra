<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import api from "../services/api";

const props = defineProps<{
  date: string;
}>();

const hours = Array.from({ length: 24 }, (_, i) => `${i}:00`);
const resources = ref<any[]>([]);
const loading = ref(false);

const fetchData = async () => {
  loading.value = true;
  try {
    const [roomsRes, bookingsRes] = await Promise.all([
      api.get("/rooms"),
      api.get("/bookings/all"),
    ]);

    resources.value = roomsRes.data.map((room: any) => {
      const roomBookings = bookingsRes.data
        .filter((b: any) => {
          // Filter by room and selected date
          const bDate = new Date(b.start_time).toISOString().split("T")[0];
          return b.room_id === room.id && bDate === props.date;
        })
        .map((b: any) => {
          const start = new Date(b.start_time);
          const end = new Date(b.end_time);

          const startHour = start.getHours() + start.getMinutes() / 60;
          const endHour = end.getHours() + end.getMinutes() / 60;

          return {
            start: startHour,
            end: endHour,
            color:
              b.status === "approved"
                ? "bg-accent"
                : b.status === "rejected"
                  ? "border-2 border-red-500 bg-red-500/10"
                  : "bg-surface border border-border",
            label: `AL-${b.id}`,
            conflict: b.status === "rejected",
          };
        });

      return {
        name: room.name,
        allocations: roomBookings,
      };
    });
  } catch (err) {
    console.error("Timeline data fetch error");
  } finally {
    loading.value = false;
  }
};

watch(
  () => props.date,
  () => {
    fetchData();
  },
);

onMounted(() => {
  fetchData();
});
</script>

<template>
  <div class="flex flex-col h-full bg-background select-none">
    <!-- Timeline Header -->
    <div class="flex border-b border-border bg-surface sticky top-0 z-10">
      <div
        class="w-48 border-r border-border p-4 text-[10px] font-bold text-muted uppercase tracking-widest leading-none flex items-center"
      >
        RESOURCE_ID
      </div>
      <div class="flex-1 flex overflow-x-auto no-scrollbar">
        <div
          v-for="h in hours"
          :key="h"
          class="flex-none w-20 border-r border-border/50 p-3 text-[10px] font-mono text-muted text-center"
        >
          {{ h }}
        </div>
      </div>
    </div>

    <!-- Timeline Body -->
    <div class="flex-1 overflow-y-auto overflow-x-hidden">
      <div
        v-for="res in resources"
        :key="res.name"
        class="flex border-b border-border group hover:bg-surface/50 transition-colors"
      >
        <!-- Row Header -->
        <div
          class="w-48 border-r border-border p-4 font-mono text-xs text-primary bg-surface/30 flex items-center"
        >
          {{ res.name }}
        </div>

        <!-- Row Cells/Grid -->
        <div class="flex-1 relative h-16 flex">
          <!-- Background Grid Lines -->
          <div
            v-for="i in 24"
            :key="i"
            class="flex-none w-20 border-r border-border/30 h-full"
          ></div>

          <!-- Allocation Blocks -->
          <div
            v-for="(alloc, idx) in res.allocations"
            :key="idx"
            class="absolute top-2 bottom-2 rounded-sm flex items-center justify-center p-2 overflow-hidden shadow-lg transition-transform hover:scale-[1.02] cursor-pointer"
            :class="alloc.color"
            :style="{
              left: alloc.start * 80 + 'px',
              width: (alloc.end - alloc.start) * 80 + 'px',
            }"
          >
            <span
              class="text-[9px] font-bold font-mono tracking-tighter truncate"
              :class="alloc.conflict ? 'text-red-500' : 'text-white'"
            >
              {{ alloc.label }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Legend -->
    <div
      class="p-4 border-t border-border bg-surface flex items-center gap-6 text-[10px] font-mono text-muted uppercase"
    >
      <div class="flex items-center gap-2">
        <span class="w-2 h-2 bg-accent rounded-sm"></span> ALLOCATED
      </div>
      <div class="flex items-center gap-2">
        <span
          class="w-2 h-2 border border-red-500 bg-red-500/20 rounded-sm"
        ></span>
        CONFLICT
      </div>
      <div class="flex items-center gap-2">
        <span class="w-2 h-2 bg-neutral-800 rounded-sm"></span> RESERVED
      </div>
    </div>
  </div>
</template>

<style scoped>
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
