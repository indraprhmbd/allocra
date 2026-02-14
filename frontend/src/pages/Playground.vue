<script setup lang="ts">
import { ref, onMounted } from "vue";
import { IconFlask, IconPlayerPlay, IconDice } from "@tabler/icons-vue";
import api from "../services/api";
import { useToast } from "../composables/useToast";

const { toast } = useToast();
const loading = ref(false);
const resources = ref<any[]>([]);

const fetchResources = async () => {
  try {
    const res = await api.get("/rooms");
    resources.value = res.data;
  } catch (err) {
    console.error("Failed to fetch resources");
  }
};

const handleSimulate = async () => {
  if (resources.value.length === 0) {
    toast("No resources available for simulation", "error");
    return;
  }

  loading.value = true;
  try {
    // Pick a random resource
    const randomRoom =
      resources.value[Math.floor(Math.random() * resources.value.length)];

    // Generate a future time window (1 hour from now, 2 hour duration)
    // Adding 5 min buffer just to be safe with the backend's grace period
    const start = new Date();
    start.setMinutes(start.getMinutes() + 5);

    const end = new Date(start);
    end.setHours(end.getHours() + 2);

    const payload = {
      room_id: randomRoom.id,
      user_id: 1, // System default
      start_time: start.toISOString(),
      end_time: end.toISOString(),
    };

    await api.post("/bookings", payload);
    toast(
      `Simulation Successful: Node ${randomRoom.name} requested`,
      "success",
    );
  } catch (err: any) {
    const errorMsg = err.response?.data?.error || "Simulation failed";
    toast(errorMsg, "error");
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchResources();
});
</script>

<template>
  <div class="max-w-4xl mx-auto space-y-12 py-8">
    <div class="text-center space-y-4">
      <div class="flex justify-center">
        <div class="p-4 bg-accent/10 rounded-full border border-accent/20">
          <IconFlask :size="48" class="text-accent" />
        </div>
      </div>
      <h2 class="text-3xl font-bold text-primary tracking-tighter">
        ALLOCRA_PLAYGROUND
      </h2>
      <p class="text-muted text-sm max-w-md mx-auto">
        Bypass manual data entry. Use the deterministic simulation engine to
        stress-test resource availability and preemption logic.
      </p>
    </div>

    <div
      class="bg-surface border border-border p-12 rounded-sm text-center space-y-8 shadow-2xl shadow-accent/5"
    >
      <div class="space-y-2">
        <h3 class="text-xs font-bold text-accent uppercase tracking-[0.2em]">
          Automated_Verification_Engine
        </h3>
        <p class="text-sm text-muted">
          Click the button below to generate a random high-validity allocation
          request.
        </p>
      </div>

      <button
        @click="handleSimulate"
        :disabled="loading"
        class="group relative inline-flex items-center gap-3 bg-accent hover:bg-accent-hover text-white text-sm font-bold uppercase tracking-widest px-12 py-5 rounded-sm transition-all shadow-[0_0_30px_rgba(var(--color-accent),0.3)] disabled:opacity-50 disabled:cursor-not-allowed overflow-hidden"
      >
        <div
          class="absolute inset-0 bg-white/10 translate-y-full group-hover:translate-y-0 transition-transform duration-300"
        ></div>
        <IconPlayerPlay v-if="!loading" :size="20" class="relative z-10" />
        <IconDice v-else :size="20" class="relative z-10 animate-spin" />
        <span class="relative z-10">{{
          loading ? "EXECUTING_SIMULATION..." : "EXECUTE_SIMULATION"
        }}</span>
      </button>

      <div class="flex justify-center gap-12 pt-8">
        <div class="text-center">
          <div class="text-lg font-mono text-primary font-bold">1-CLICK</div>
          <div
            class="text-[10px] text-muted uppercase tracking-wider font-medium"
          >
            Efficiency
          </div>
        </div>
        <div class="w-px h-10 bg-border"></div>
        <div class="text-center">
          <div class="text-lg font-mono text-primary font-bold">STOCHASTIC</div>
          <div
            class="text-[10px] text-muted uppercase tracking-wider font-medium"
          >
            Logic
          </div>
        </div>
        <div class="w-px h-10 bg-border"></div>
        <div class="text-center">
          <div class="text-lg font-mono text-primary font-bold">24/7</div>
          <div
            class="text-[10px] text-muted uppercase tracking-wider font-medium"
          >
            Availability
          </div>
        </div>
      </div>
    </div>

    <!-- Technical Advisory -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="p-6 bg-surface/50 border border-border rounded-sm space-y-3">
        <div class="flex items-center gap-2 text-accent">
          <IconFlask :size="16" />
          <span class="text-[10px] font-bold uppercase tracking-widest"
            >Simulation_Parameters</span
          >
        </div>
        <ul class="text-[11px] text-muted font-mono space-y-1">
          <li>> TARGET: RANDOM_ONLINE_NODE</li>
          <li>> SCOPE: T + 5MIN (GLOBAL_SYNC)</li>
          <li>> DURATION: 120min (FIXED)</li>
        </ul>
      </div>
      <div class="p-6 bg-surface/50 border border-border rounded-sm space-y-3">
        <div class="flex items-center gap-2 text-yellow-500">
          <IconPlayerPlay :size="16" />
          <span class="text-[10px] font-bold uppercase tracking-widest"
            >Logic_Bypass</span
          >
        </div>
        <p class="text-[11px] text-muted leading-relaxed">
          THIS ACTION BYPASSES MANUAL TIME SELECTION. SIMULATION REQUESTS ARE
          SUBJECT TO STANDARD CONFLICT DETECTION (idx_bookings_room_time).
        </p>
      </div>
    </div>
  </div>
</template>
