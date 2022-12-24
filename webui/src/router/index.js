import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginForm from "@/views/Login.vue";
import ManageProfile from "@/views/ManageProfile.vue";
import Upload from "@/views/Upload.vue";
import Search from "@/views/Search.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/login', component: LoginForm},
		{path: '/users/:user', component: ManageProfile},
		{path: '/upload', component: Upload},
		{path: '/search', component: Search},
		{path: '/users/me', component: ManageProfile}

	]
})

export default router
