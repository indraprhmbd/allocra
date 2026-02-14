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

    if (systemRes.data) {
      if (stats.value[0])
        stats.value[0].value = (systemRes.data.total_bookings || 0).toString();
      if (stats.value[1])
        stats.value[1].value = (systemRes.data.active_bookings || 0).toString();
      if (stats.value[3])
        stats.value[3].value = (systemRes.data.conflicts || 0).toString();

      const util = systemRes.data.utilization || 0;
      systemStats.value = {
        ...systemStats.value,
        cpu_usage: Math.round(util),
        memory_usage: Math.round(util * 0.8 + 10), // Simulated relation
        io_wait: Math.round(util * 0.3),
        status: util > 90 ? "OVERLOAD" : util > 0 ? "ONLINE" : "IDLE",
      };
    }

    if (stats.value[0]) stats.value[0].value = roomsRes.data.length.toString();
    if (stats.value[2])
      stats.value[2].value = bookingsRes.data
        .filter((b: any) => b.status === "pending")
        .length.toString();

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
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold tracking-tight text-white">Dashboard</h1>
      <div class="flex items-center gap-2">
        <span class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
        <span class="text-xs font-mono text-accent">LIVE_STREAM</span>
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6">
      <StatCard
        v-for="stat in stats"
        :key="stat.label"
        v-bind="stat"
        class="bg-surface border border-border p-4"
      />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Activity Feed -->
      <div class="lg:col-span-2 space-y-4">
        <h2 class="text-lg font-semibold text-white flex items-center gap-2">
          <IconActivity :size="20" class="text-accent" />
          Latest Allocation Events
        </h2>
        <div class="bg-surface border border-border rounded-sm overflow-hidden">
          <div class="overflow-x-auto">
            <DataTable
              :columns="columns"
              :items="latestEvents"
              :loading="false"
            >
              <template #status="{ value }">
                <StatusBadge :status="value" />
              </template>
            </DataTable>
          </div>
        </div>
      </div>

      <!-- System Health -->
      <div class="space-y-4 hidden lg:block">
        <h2 class="text-lg font-semibold text-white flex items-center gap-2">
          <IconCpu :size="20" class="text-accent" />
          Engine Load
        </h2>

        <div class="bg-surface border border-border rounded-sm p-5">
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
