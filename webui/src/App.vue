<script setup>
import { RouterLink, RouterView } from 'vue-router'

</script>
<script>
export default {
  data: function () {
    return {
      username: null,
      loggedIn: false
    }
  },
  methods:{
    logout(){
      sessionStorage.removeItem("username")
      sessionStorage.removeItem("token")
      this.loggedIn = false
      this.$router.push("/login")
      alert("Logged out")
    },
    logging(){
      this.loggedIn = true
      this.username = sessionStorage.getItem("username")
    }

  },
}
</script>

<template>

	<header class="navbar sticky-top flex-md-nowrap p-0 shadow" style="background-color: #1565c0">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" style="background-color: #003c8f; color: white" href="#/">WASA Photo</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
    <button class="btn btn-dark me-2" v-if="loggedIn" @click="logout()" >Logout</button>
  </header>

	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>General</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item" v-if="loggedIn">
							<RouterLink to="/" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
								Home
							</RouterLink>
						</li>
						<li class="nav-item" v-if="!loggedIn">
							<RouterLink to="/login" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
								Login
							</RouterLink>
						</li>
            <li class="nav-item" v-if="loggedIn">
              <RouterLink :to="`/users/${this.username}`" class="nav-link">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#file-text"/></svg>
                My account
              </RouterLink>
            </li>
            <li class="nav-item" v-if="loggedIn">
              <RouterLink to="/upload" class="nav-link">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
                Upload a photo
              </RouterLink>
            </li>
          </ul>
				</div>
			</nav>

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView :key="$route.fullPath" @logging="logging()" /> <!-- The key forces vue to reaload every component when the URL changes -->
			</main>
		</div>
	</div>
</template>

<style>
</style>
