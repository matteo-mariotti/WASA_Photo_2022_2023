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
      username: sessionStorage.getItem('username'),
      isOwner: null,
      page: 0,
      isOpen: ref(false),
      newUsername: null,
      Infomsg: null,
      moreButton: null
    }
  },
  methods: {
    async loadContent() {
      try {
        this.ready = false
        let response = await backend.get(`/users/${this.$route.params.user}`);
        this.handleResponse(response);
        this.ready = true
      } catch (error) {
        this.handleError(error);
      }
    },
    handleResponse(response) {
      this.errormsg = null
      this.userData = response.data

      this.isOwner = this.userData.username === sessionStorage.getItem("username")
      this.page = this.page + 1
      this.moreButton = true
      console.log(response)
    },
    handleError(error) {

      if (error.response.status == 401) {
        this.errormsg = "You need to login first"
      } else if (error.response.status == 403) {
        this.errormsg = "User ID not found"
      } else {
        // Print the error from the server (for debugging)
        this.errormsg = error.response.data.message
      }
      // Log the error on console
      console.log(error)
    },
    async loadMorePhotos(offset) {
      // TODO Fare la chiamata per ottenere le foto a partire da un certo offset
      try {
        let response = await backend.get(`/users/${this.$route.params.user}?page=${this.page}`);
        console.log(response)
        if (response.status == 200) {
          response.data.photos.forEach(photo => {
            this.userData.photos.push(photo)
          })
        }
        this.page = this.page + 1
        if (response.status == 204){
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
        setTimeout(this.reloadPage, 1000);
      } catch (error) {
        this.Modalerrormsg = error.response.data.message
        console.log(error)
      }
    },
    reloadPage() {
      location.reload();
    }
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
    <!-- Profile stats -->
    <div class="d-flex justify-content-around p-2"
         style="font-size: 2vw;background-color: rgba(137,137,137,0.2);border-radius: 2vw;" v-if="ready">
      <div v-if="ready">Photo: {{ userData.photoNumber }}</div>
      <div v-if="ready">Follower: {{ userData.follower }}</div>
      <div v-if="ready">Following: {{ userData.following }}</div>
    </div>

    <!-- Images -->
    <div class="d-grid justify-content-center" id="imageList">
      <Card v-bind:imageData="image" v-if="ready" v-for="image in userData.photos"
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
