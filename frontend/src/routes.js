import Controller from "./pages/Controller.vue";
import Home from "./pages/Home.vue";
import Viewer from "./pages/Viewer.vue";

export default [
  { path: "/", component: Home },
  { path: "/controller", component: Controller },
  { path: "/viewer", component: Viewer }
];
