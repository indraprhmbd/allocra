export interface Resource {
  id: number;
  name: string;
  capacity: number;
  type: "exclusive" | "shared";
  usage: number;
  status: "online" | "offline" | "maintenance";
  created_at: string;
}

export interface AllocationRequest {
  id: number;
  resource_id: number;
  user_id: number;
  start_time: string;
  end_time: string;
  capacity_required: number;
  priority: "low" | "medium" | "high" | "critical";
  status: "pending" | "allocated" | "rejected";
  created_at: string;
}

export interface AllocationResult {
  success: boolean;
  allocation_id?: number;
  conflicts?: AllocationRequest[];
  message?: string;
}

export interface Stat {
  label: string;
  value: string | number;
  trend?: number;
  icon?: string;
}
