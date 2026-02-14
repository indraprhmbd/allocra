<script setup lang="ts">
import { ref, onMounted } from "vue";
import {
  IconFlask,
  IconPlayerPlay,
  IconDice,
  IconTerminal2,
  IconTrash,
  IconAlertCircle,
} from "@tabler/icons-vue";
import api from "../services/api";

interface TraceLog {
  id: string;
  timestamp: string;
  type: "success" | "conflict" | "error" | "info";
  message: string;
  details?: string;
}

const loading = ref(false);
const resources = ref<any[]>([]);
const batchSize = ref(5);
const mode = ref<"sequential" | "parallel">("sequential");
const logs = ref<TraceLog[]>([]);

const fetchResources = async () => {
  try {
    const res = await api.get("/rooms");
    resources.value = res.data;
  } catch (err) {
    addLog("error", "SYSTEM", "Failed to fetch resource nodes");
  }
};

const addLog = (
  type: TraceLog["type"],
  tag: string,
  message: string,
  details?: string,
) => {
  const log: TraceLog = {
    id: Math.random().toString(36).substring(7),
    timestamp: new Date().toLocaleTimeString(),
    type,
    message: `[${tag}] ${message}`,
    details,
  };
  logs.value.unshift(log); // Newest at top
  if (logs.value.length > 50) logs.value.pop();
  localStorage.setItem("playground_logs", JSON.stringify(logs.value));
};

const clearLogs = () => {
  logs.value = [];
  localStorage.removeItem("playground_logs");
};

const handleReset = async () => {
  if (!confirm("WARNING: This will delete ALL allocations. Continue?")) return;

  loading.value = true;
  try {
    await api.post("/allocations/reset");
    addLog("info", "SYSTEM", "Database reset successful");
  } catch (err: any) {
    addLog("error", "SYSTEM", "Reset failed", err.message);
  } finally {
    loading.value = false;
  }
};

const simulateOne = async (index: number) => {
  const randomRoom =
    resources.value[Math.floor(Math.random() * resources.value.length)];

  // Random staggered start times in the next 24 hours
  const start = new Date();
  const staggeredMinutes = 5 + index * 15 + Math.floor(Math.random() * 60);
  start.setMinutes(start.getMinutes() + staggeredMinutes);

  const end = new Date(start);
  end.setHours(end.getHours() + 1 + Math.floor(Math.random() * 2));

  try {
    const payload = {
      room_id: randomRoom.id,
      user_id: 1,
      start_time: start.toISOString(),
      end_time: end.toISOString(),
    };

    await api.post("/bookings", payload);
    addLog(
      "success",
      randomRoom.name,
      `Allocation approved for ${start.toLocaleTimeString()} - ${end.toLocaleTimeString()}`,
    );
  } catch (err: any) {
    const errorMsg = err.response?.data?.error || "Unknown error";
    const status = err.response?.status;

    if (status === 409 || errorMsg.includes("conflict")) {
      addLog(
        "conflict",
        randomRoom.name,
        `Resource conflict detected`,
        errorMsg,
      );
    } else {
      addLog("error", randomRoom.name, `Failed: ${errorMsg}`);
    }
  }
};

const handleFloodSimulation = async () => {
  if (resources.value.length === 0) {
    addLog("error", "SYSTEM", "Aborting: No resources found");
    return;
  }

  loading.value = true;
  addLog(
    "info",
    "FLOOD",
    `Initiating ${mode.value} batch of ${batchSize.value} requests...`,
  );

  if (mode.value === "sequential") {
    // Execute sequentially
    for (let i = 0; i < batchSize.value; i++) {
      await simulateOne(i);
      await new Promise((r) => setTimeout(r, 100));
    }
  } else {
    // Execute in parallel
    const promises = Array.from({ length: batchSize.value }, (_, i) =>
      simulateOne(i),
    );
    await Promise.all(promises);
  }

  addLog("info", "FLOOD", `Batch execution completed.`);
  loading.value = false;
};

onMounted(() => {
  fetchResources();
  addLog("info", "SYSTEM", "Playground engine initialized. REPL ready.");
});
</script>

<template>
  <div class="max-w-6xl mx-auto py-8 px-4 space-y-8">
    <!-- Header -->
    <div class="flex items-start justify-between">
      <div class="space-y-1">
        <h2
          class="text-3xl font-bold text-primary tracking-tighter flex items-center gap-3"
        >
          <IconFlask class="text-accent" /> ALLOCRA_PLAYGROUND
        </h2>
        <p
          class="text-muted text-sm font-mono uppercase tracking-widest bg-accent/5 px-2 py-0.5 rounded-sm inline-block"
        >
          Deterministic_Stress_Testing_Unit_v2
        </p>
      </div>
      <div class="flex items-center gap-4">
        <div class="text-right">
          <p class="text-[10px] font-mono text-muted uppercase">
            Engine_Status
          </p>
          <p class="text-green-500 font-bold font-mono">NOMINAL</p>
        </div>
        <div
          class="w-10 h-10 rounded-sm bg-accent/10 border border-accent/20 flex items-center justify-center"
        >
          <IconTerminal2 class="text-accent" />
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
      <!-- Controls Panel -->
      <div class="lg:col-span-4 space-y-6">
        <div class="bg-surface border border-border p-6 rounded-sm space-y-6">
          <h3
            class="text-xs font-bold text-primary uppercase tracking-widest flex items-center gap-2"
          >
            <IconDice :size="14" /> Simulation_Parameters
          </h3>

          <div class="space-y-4">
            <div class="space-y-2">
              <div
                class="flex justify-between text-[10px] font-mono text-muted uppercase"
              >
                <span>Request_Flood_Level</span>
                <span class="text-accent">{{ batchSize }} REQS</span>
              </div>
              <input
                type="range"
                v-model.number="batchSize"
                min="1"
                max="50"
                class="w-full accent-accent bg-background h-1.5 rounded-full appearance-none border border-border"
              />
            </div>

            <div
              class="p-4 bg-background/50 border border-dashed border-border rounded-sm space-y-2"
            >
              <p class="text-[10px] text-muted font-mono leading-relaxed">
                > TARGET: {{ resources.length }} NODES ACTIVE<br />
                > STRATEGY: STOCHASTIC_TIME_SLOTS<br />
                > MODE: {{ mode.toUpperCase() }}_TRACE
              </p>
            </div>

            <!-- Mode Toggle -->
            <div
              class="flex items-center gap-2 bg-background p-1 border border-border rounded-sm"
            >
              <button
                @click="mode = 'sequential'"
                class="flex-1 py-1.5 text-[10px] font-bold uppercase tracking-wider rounded-sm transition-colors"
                :class="
                  mode === 'sequential'
                    ? 'bg-accent text-white shadow-sm'
                    : 'text-muted hover:text-primary'
                "
              >
                Sequential
              </button>
              <button
                @click="mode = 'parallel'"
                class="flex-1 py-1.5 text-[10px] font-bold uppercase tracking-wider rounded-sm transition-colors"
                :class="
                  mode === 'parallel'
                    ? 'bg-accent text-white shadow-sm'
                    : 'text-muted hover:text-primary'
                "
              >
                Parallel
              </button>
            </div>
          </div>

          <button
            @click="handleFloodSimulation"
            :disabled="loading"
            class="w-full relative group bg-accent hover:bg-accent-hover text-white text-xs font-bold uppercase tracking-[0.2em] py-4 rounded-sm transition-all disabled:opacity-50 overflow-hidden"
          >
            <div
              class="absolute inset-0 bg-white/10 -translate-x-full group-hover:translate-x-0 transition-transform duration-500"
            ></div>
            <span class="relative z-10 flex items-center justify-center gap-2">
              <IconPlayerPlay v-if="!loading" :size="16" />
              <IconDice v-else :size="16" class="animate-spin" />
              {{ loading ? "EXECUTING_STORM..." : "INITIATE_FLOOD" }}
            </span>
          </button>
        </div>

        <div
          class="p-4 bg-yellow-500/5 border border-yellow-500/20 rounded-sm flex gap-3"
        >
          <IconAlertCircle :size="16" class="text-yellow-500 shrink-0 mt-0.5" />
          <p class="text-[10px] text-muted leading-relaxed">
            CAUTION: FLOOD_MODE generates real database entries. High volume may
            affect timeline visualization performance.
          </p>
        </div>

        <button
          @click="handleReset"
          class="w-full border border-red-500/30 hover:bg-red-500/10 text-red-500 text-xs font-bold uppercase tracking-widest py-3 rounded-sm transition-colors flex items-center justify-center gap-2"
        >
          <IconTrash :size="14" /> PURGE_DATABASE
        </button>
      </div>

      <!-- Live Tracer Panel -->
      <div class="lg:col-span-8 flex flex-col h-[600px]">
        <div
          class="bg-background border border-border rounded-sm flex flex-col h-full shadow-inner"
        >
          <!-- Tracer Header -->
          <div
            class="p-3 border-b border-border bg-surface/50 flex items-center justify-between"
          >
            <div class="flex items-center gap-2">
              <div class="flex gap-1">
                <div class="w-2.5 h-2.5 rounded-full bg-red-500/50"></div>
                <div class="w-2.5 h-2.5 rounded-full bg-yellow-500/50"></div>
                <div class="w-2.5 h-2.5 rounded-full bg-green-500/50"></div>
              </div>
              <span
                class="text-[10px] font-mono font-bold text-muted ml-4 uppercase tracking-widest"
                >Live_Engine_Tracer</span
              >
            </div>
            <button
              @click="clearLogs"
              class="text-muted hover:text-red-400 transition-colors"
              title="Clear Logs"
            >
              <IconTrash :size="14" />
            </button>
          </div>

          <!-- Tracer Log Area -->
          <div
            class="flex-1 overflow-y-auto p-4 font-mono text-[11px] space-y-2 selection:bg-accent/30"
          >
            <div
              v-if="logs.length === 0"
              class="h-full flex items-center justify-center text-muted/30 uppercase tracking-widest italic"
            >
              Awaiting_Input...
            </div>
            <TransitionGroup name="log">
              <div
                v-for="log in logs"
                :key="log.id"
                class="flex gap-3 group animate-in slide-in-from-left-1 duration-200"
              >
                <span class="text-muted/50 shrink-0">{{ log.timestamp }}</span>
                <div class="space-y-1">
                  <div
                    :class="{
                      'text-green-400': log.type === 'success',
                      'text-red-400': log.type === 'conflict',
                      'text-accent': log.type === 'info',
                      'text-red-500 font-bold': log.type === 'error',
                    }"
                  >
                    <span v-if="log.type === 'success'" class="mr-1">✓</span>
                    <span v-if="log.type === 'conflict'" class="mr-1">⚠</span>
                    <span v-if="log.type === 'error'" class="mr-1">❌</span>
                    {{ log.message }}
                  </div>
                  <div
                    v-if="log.details"
                    class="text-muted/60 pl-4 border-l border-border/50 text-[10px] leading-relaxed italic"
                  >
                    {{ log.details }}
                  </div>
                </div>
              </div>
            </TransitionGroup>
          </div>

          <!-- Tracer Footer -->
          <div
            class="p-2 border-t border-border bg-surface/30 px-4 flex justify-between items-center"
          >
            <div class="flex gap-4">
              <span class="text-[9px] text-muted uppercase"
                >Trace_Buffer: {{ logs.length }}/50</span
              >
            </div>
            <div class="flex items-center gap-1.5 animate-pulse">
              <div class="w-1.5 h-1.5 rounded-full bg-accent"></div>
              <span
                class="text-[9px] text-accent uppercase font-bold tracking-tighter"
                >Engine_Listening</span
              >
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.log-enter-active {
  transition: all 0.3s ease-out;
}
.log-enter-from {
  opacity: 0;
  transform: translateX(-10px);
}

/* Custom scrollbar for terminal feel */
::-webkit-scrollbar {
  width: 4px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  background: rgba(var(--color-accent-rgb), 0.2);
  border-radius: 2px;
}
::-webkit-scrollbar-thumb:hover {
  background: rgba(var(--color-accent-rgb), 0.4);
}
</style>
