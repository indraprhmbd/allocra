<script setup lang="ts">
import { ref, onMounted } from "vue";
import {
  IconDatabase,
  IconArrowsExchange,
  IconAlertTriangle,
  IconActivity,
  IconClock,
} from "@tabler/icons-vue";
import StatCard from "../components/StatCard.vue";
import DataTable, { type Column } from "../components/DataTable.vue";
import api from "../services/api";

const stats = ref([
  { label: "Total Resources", value: "0", icon: IconDatabase, trend: 0 },
  {
    label: "Active Allocations",
    value: "0",
    icon: IconArrowsExchange,
    trend: 0,
  },
  { label: "Pending Requests", value: "0", icon: IconClock, trend: 0 },
  {
    label: "Detected Conflicts",
    value: "0",
    icon: IconAlertTriangle,
    trend: 0,
  },
]);

const systemStats = ref({
  cpu_usage: 0,
  memory_usage: 0,
  io_wait: 0,
  status: "OFFLINE",
});

const fetchDashboardData = async () => {
  try {
    const [roomsRes, bookingsRes, systemRes] = await Promise.all([
      api.get("/rooms"),
      api.get("/bookings/all"),
      api.get("/system/stats"),
    ]);

    // Map backend stats to UI
    if (systemRes.data) {
      // Update top cards
      if (stats.value[0])
        stats.value[0].value = (systemRes.data.total_bookings || 0).toString();
      if (stats.value[1])
        stats.value[1].value = (systemRes.data.active_bookings || 0).toString();

      // Calculate pending (total - active - rejected) or fetch if available.
      // For now, let's stick to frontend filter for specific statuses not in SystemStats or add pending to SystemStats later.
      // Actually, let's keep the filter for properties NOT in SystemStats to be safe, or mix/match.
      // SystemStats has Total, Active, Conflicts.

      // Let's use the robust backend numbers where possible
      if (stats.value[1])
        stats.value[1].value = (systemRes.data.active_bookings || 0).toString();
      if (stats.value[3])
        stats.value[3].value = (systemRes.data.conflicts || 0).toString();

      // Simulate Engine Load based on Utilization
      const util = systemRes.data.utilization || 0;
      systemStats.value = {
        ...systemStats.value,
        cpu_usage: Math.round(util),
        memory_usage: Math.round(util * 0.8 + 10), // Simulated relation
        io_wait: Math.round(util * 0.3),
        status: util > 90 ? "OVERLOAD" : util > 0 ? "ONLINE" : "IDLE",
      };
    }

    // Fallback/Supplemental from bookings list for things like "Total Resources" (rooms count)
    if (stats.value[0]) stats.value[0].value = roomsRes.data.length.toString();
    if (stats.value[2])
      stats.value[2].value = bookingsRes.data
        .filter((b: any) => b.status === "pending")
        .length.toString();

    // Map latest events from real bookings
    latestEvents.value = bookingsRes.data
      .sort(
        (a: any, b: any) =>
          new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
      )
      .slice(0, 5)
      .map((b: any) => ({
        timestamp: new Date(b.created_at).toLocaleString(),
        resource: b.room_name || `NODE-ID-${b.room_id}`,
        event:
          b.status === "approved"
            ? "ALLOCATION_START"
            : b.status === "pending"
              ? "RESERVATION_REQUEST"
              : "CONFLICT_DETECTED",
        status: b.status.toUpperCase(),
      }));
  } catch (err) {
    console.error("Dashboard sync error");
  }
};

onMounted(() => {
  fetchDashboardData();
});

const columns: Column[] = [
  { key: "timestamp", label: "Timestamp", class: "w-48 font-mono text-xs" },
  { key: "resource", label: "Resource" },
  { key: "event", label: "Event Type" },
  { key: "status", label: "Status", class: "text-right" },
];

const latestEvents = ref<any[]>([]);
</script>

<template>
  <div class="space-y-8">
    <div class="flex items-end justify-between">
      <div>
        <h2 class="text-2xl font-bold text-primary tracking-tight">
          System Dashboard
        </h2>
        <p class="text-muted text-sm mt-1">
          Real-time resource utilization and allocation metrics.
        </p>
      </div>
      <div
        class="flex items-center gap-2 text-xs font-mono text-muted uppercase tracking-widest border-b border-border pb-1"
      >
        <IconActivity :size="14" class="text-accent animate-pulse" />
        LIVE_STREAM_CONNECTED: {{ new Date().toLocaleTimeString() }}
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <StatCard v-for="stat in stats" :key="stat.label" v-bind="stat" />
    </div>

    <!-- Main Content Area -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <div
        class="lg:col-span-2 bg-surface border border-border rounded-sm overflow-hidden flex flex-col"
      >
        <div
          class="p-5 border-b border-border flex items-center justify-between"
        >
          <h3
            class="text-sm font-bold uppercase tracking-widest text-primary flex items-center gap-2"
          >
            <IconArrowsExchange :size="16" class="text-accent" />
            Latest Allocation Events
          </h3>
          <button
            class="text-[10px] font-mono text-muted hover:text-accent transition-colors"
          >
            VIEW_ALL_LOGS
          </button>
        </div>
        <DataTable :columns="columns" :items="latestEvents">
          <template #status="{ value }">
            <span
              :class="[
                value === 'SUCCESS' || value === 'OK'
                  ? 'text-green-500'
                  : value === 'PENDING'
                    ? 'text-accent'
                    : 'text-red-500',
              ]"
              class="text-[10px] font-mono px-1.5 py-0.5 border border-current/20 bg-current/5 rounded-sm"
            >
              {{ value }}
            </span>
          </template>
        </DataTable>
      </div>

      <div class="space-y-8">
        <div class="bg-surface border border-border rounded-sm p-5">
          <h3
            class="text-sm font-bold uppercase tracking-widest text-primary mb-4 flex items-center gap-2"
          >
            <IconActivity :size="16" class="text-accent" />
            Engine Load
          </h3>
          <div class="space-y-4">
            <div>
              <div
                class="flex justify-between text-[10px] font-mono text-muted mb-1.5 uppercase"
              >
                <span>CPU Usage</span>
                <span>{{ systemStats.cpu_usage }}%</span>
              </div>
              <div class="h-1 bg-background rounded-full overflow-hidden">
                <div
                  class="h-full bg-accent transition-all duration-500"
                  :style="{ width: systemStats.cpu_usage + '%' }"
                ></div>
              </div>
            </div>
            <div>
              <div
                class="flex justify-between text-[10px] font-mono text-muted mb-1.5 uppercase"
              >
                <span>Memory</span>
                <span>{{ systemStats.memory_usage }}%</span>
              </div>
              <div class="h-1 bg-background rounded-full overflow-hidden">
                <div
                  class="h-full bg-accent transition-all duration-500"
                  :style="{ width: systemStats.memory_usage + '%' }"
                ></div>
              </div>
            </div>
            <div>
              <div
                class="flex justify-between text-[10px] font-mono text-muted mb-1.5 uppercase"
              >
                <span>I/O Wait</span>
                <span>{{ systemStats.io_wait }}%</span>
              </div>
              <div class="h-1 bg-background rounded-full overflow-hidden">
                <div
                  class="h-full bg-accent transition-all duration-500"
                  :style="{ width: systemStats.io_wait + '%' }"
                ></div>
              </div>
            </div>
          </div>
        </div>

        <div
          class="bg-surface border border-border rounded-sm p-5 border-l-4 border-l-accent"
        >
          <h3
            class="text-xs font-bold uppercase tracking-widest text-accent mb-2"
          >
            Notice
          </h3>
          <p class="text-xs text-muted leading-relaxed">
            Automatic conflict resolution is currently set to
            <span class="text-primary">PRESERVE_EXISTING</span>. Change this in
            system settings if preemption is required.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
