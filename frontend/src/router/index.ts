import { createRouter, createWebHistory } from "vue-router";
import AppLayout from "../components/AppLayout.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      component: AppLayout,
      children: [
        {
          path: "",
          name: "Dashboard",
          component: () => import("../pages/Dashboard.vue"),
        },
        {
          path: "resources",
          name: "Resources",
          component: () => import("../pages/Resources.vue"),
        },
        {
          path: "allocations",
          name: "Allocations",
          component: () => import("../pages/Allocations.vue"),
        },
        {
          path: "timeline",
          name: "Timeline",
          component: () => import("../pages/Timeline.vue"),
        },
      ],
    },
  ],
});

export default router;
