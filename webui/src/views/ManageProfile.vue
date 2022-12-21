<script>
import backend from "@/services/axios";
import ErrorMsg from "@/components/ErrorMsg.vue";
import SuccessMsg from "@/components/SuccessMsg.vue";
import Card from "@/components/Card.vue";
import ModalV2 from "@/components/ModalV2.vue";
import {ref} from "vue";
import InfoMsg from "@/components/InfoMsg.vue";

export default {
  components: {InfoMsg, ModalV2, Card, SuccessMsg, ErrorMsg},
  data: function () {
    return {
      errormsg: null,
      successmsg: null,
      Modalsuccessmsg: null,
      Modalerrormsg: null,
      ready: false,
      userData: null,
      photos: [],
      username: sessionStorage.getItem('username'),
      isOwner: null,
      page: 0,
      isOpen: ref(false),
      newUsername: null,
      Infomsg: null,
      moreButton: null,
      banned: false,
      following: null,
    }
  },
  methods: {
    async loadContent() {
      try {
        this.ready = false
        let response = await backend.get(`/users/${sessionStorage.getItem("username")}/bans/${this.$route.params.user}`);
        if (response.status === 200) {
          if (response.data.status === true) {
            console.log("Sono qui")
            this.banned = true
            this.errormsg = "You banned this user"
          }
        }

        if (this.banned === false) {
          response = await backend.get(`/users/${this.$route.params.user}`);
          this.handleResponse(response);
          this.ready = true
          await this.isFollowing()
          await this.loadMorePhotos(this.page)
        }
      } catch (error) {
        this.handleError(error);
      }
    },
    handleResponse(response) {
      this.errormsg = null
      this.userData = response.data

      this.isOwner = this.userData.username === sessionStorage.getItem("username")
      this.moreButton = true
      console.log(response)
    },
    handleError(error) {

      if (error.response.status === 401) {
        this.errormsg = "You need to login first"
      } else if (error.response.status === 403) {
        this.errormsg = "User ID not found"
      } else {
        // Print the error from the server (for debugging)
        this.errormsg = error.response.data.message
      }
      // Log the error on console
      console.log(error)
    },
    async loadMorePhotos(offset) {
      try {
        let response = await backend.get(`/users/${this.$route.params.user}/photos?page=${offset}`);
        console.log(response)
        if (response.status === 200) {
          response.data.forEach(photo => {
            this.photos.push(photo)
          })
        }
        this.page = this.page + 1
        if (response.status === 204) {
          this.moreButton = false;
          this.Infomsg = "No more photos are available"
        }
      } catch (error) {
        console.log(error)
      }

    },
    async changeUsername() {
      try {
        let response = await backend.put(`/users/${sessionStorage.getItem("username")}/username`, {
          username: this.newUsername
        });
        this.handleResponse(response);
        sessionStorage.setItem("username", this.newUsername)
        this.Modalerrormsg = null
        this.Modalsuccessmsg = "Username updated"
        setTimeout(this.reloadNewUser, 1000);
      } catch (error) {
        this.Modalerrormsg = error.response.data.message
        console.log(error)
      }
    },
    reloadNewUser() {
      this.$router.push(`/users/${this.newUsername}`)
      this.$emit("logging")
    },
    async isBanned() {
      try {
        let response = await backend.get(`/users/${sessionStorage.getItem("username")}/bans/${this.$route.params.user}`);
        console.log(response)
        if (response.status === 200) {
          if (response.data.status === true) {
            console.log("Sono qui")
            this.banned = true
            this.errormsg = "You banned this user"
          }
        }
      } catch (error) {
        this.Modalerrormsg = error.response.data.message
        console.log(error)
      }
    },
    async isFollowing() {
      try {
        let response = await backend.get(`/users/${this.$route.params.user}/followers/${sessionStorage.getItem("username")}`);
        console.log(response)
        if (response.status === 200) {
          if (response.data.status === true) {
            this.following = true
          }
        }
      } catch (error) {
        this.Modalerrormsg = error.response.data.message
        console.log(error)
      }
    },
    async follow() {
      try {
        let response = await backend.put(`/users/${this.$route.params.user}/followers/${sessionStorage.getItem("username")}`);
        if (response.status === 204) {
          this.following = true
          this.userData.follower = this.userData.follower + 1
        }
      } catch (error) {
        this.errormsg = error.response.data.message
        console.log(error)
      }
    },
    async unfollow() {
      try {
        let response = await backend.delete(`/users/${this.$route.params.user}/followers/${sessionStorage.getItem("username")}`);
        if (response.status === 204) {
          this.following = false
          this.userData.follower = this.userData.follower - 1
        }
      } catch (error) {
        this.errormsg = error.response.data.message
        console.log(error)
      }
    },
    async ban() {
      try {
        let response = await backend.put(`/users/${sessionStorage.getItem("username")}/bans/${this.$route.params.user}`);
        if (response.status === 204) {
          location.reload()
        }
      } catch (error) {
        this.errormsg = error.response.data.message
        console.log(error)
      }
    },
    async unban() {
      try {
        let response = await backend.delete(`/users/${sessionStorage.getItem("username")}/bans/${this.$route.params.user}`);
        if (response.status === 204) {
          location.reload()
        }
      } catch (error) {
        this.errormsg = error.response.data.message
        console.log(error)
      }
    },
  },
  mounted() {
      this.loadContent();
  }
}
</script>

<template>

  <ErrorMsg v-if="errormsg" :msg="errormsg" class="mt-3"></ErrorMsg>
  <SuccessMsg v-if="successmsg" :msg="successmsg"></SuccessMsg>

  <div class="justify-content-center align-items-center">

    <!-- Username of the user -->
    <div class="h2 d-flex justify-content-evenly mt-4" v-if="ready" style="font-size:3.5vw;">

      {{ userData.username }}

      <!-- Change username if it is the owner of the profile -->

      <button @click="isOpen = true" v-if="isOwner" class="btn btn-outline-primary btn-sm"> Change my username</button>

      <ModalV2 :open="isOpen" @close="isOpen=!isOpen">
        <SuccessMsg v-if="Modalsuccessmsg" :msg="Modalsuccessmsg"></SuccessMsg>
        <ErrorMsg v-if="Modalerrormsg" :msg="Modalerrormsg"></ErrorMsg>
        <div class="d-flex">
          <input type="text" v-model="newUsername" class="form-control" placeholder="New username">
          <button @click="changeUsername()" class="btn btn-sm btn-outline-primary">Change</button>
        </div>
      </ModalV2>

    </div>

    <!-- Follow/Unfollow and Ban/Unban buttons -->
    <div class="d-flex justify-content-evenly" v-if="!isOwner">
      <button @click="follow()" v-if="!this.following && !this.banned && this.ready" class="btn btn-success btn-sm">
        Follow
      </button>
      <button @click="unfollow()" v-if="this.following && !this.banned && this.ready" class="btn btn-warning btn-sm">
        Unfollow
      </button>
      <button @click="ban()" v-if="!this.banned && this.ready" class="btn btn-danger btn-sm">Ban</button>
      <button @click="unban()" v-if="this.banned" class="btn btn-danger btn-sm">Unban</button>
    </div>

    <hr v-if="ready">

    <!-- Profile stats -->
    <div class="d-flex justify-content-around p-1 m-1"
         style="font-size: 1.5vw;" v-if="ready">
      <div v-if="ready">Photos: {{ userData.photoNumber }}</div>
      <div v-if="ready">Followers: {{ userData.follower }}</div>
      <div v-if="ready">Followings: {{ userData.following }}</div>
    </div>

    <hr v-if="ready">
    <!-- Images -->
    <div class="d-grid justify-content-center" id="imageList">
      <Card v-bind:imageData="image" v-if="ready" v-for="image in photos"
            v-bind:username="this.username"></Card>
    </div>
    <div class="d-grid justify-content-center">
      <button @click="loadMorePhotos(this.page)" v-if="moreButton" type="button" class="btn btn-outline-dark m-3">Load
        More
      </button>
      <InfoMsg :msg="this.Infomsg" v-if="this.Infomsg"></InfoMsg>
    </div>


  </div>

</template>

<style>
</style>
