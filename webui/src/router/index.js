import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginForm from "@/views/Login.vue";
import ProfilePage from "@/views/Profile.vue"

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/login', component: LoginForm},
		{path: '/users/:userID', component: ProfilePage}
	]
})

export default router
