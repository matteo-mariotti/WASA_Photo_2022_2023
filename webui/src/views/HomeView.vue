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
      moreButton: true,
      banned: false,
      following: null,
    }
  },
  methods: {
    async loadContent() {
      try {
        this.ready = false
        await this.loadMorePhotos(this.page)
      } catch (error) {
        this.handleError(error);
      }
    },
    handleError(error) {

      if (error.response.status === 401) {
        this.errormsg = "You need to login first"
      } else {
        // Print the error from the server (for debugging)
        this.errormsg = error.response.data.message
      }
      // Log the error on console
      console.log(error)
    },
    async loadMorePhotos(offset) {
      try {
        let response = await backend.get(`/stream?page=${offset}`);
        console.log(response)
        if (response.status === 200) {
          response.data.forEach(photo => {
            this.photos.push(photo)
          })
        }
        this.page = this.page + 1
        this.ready = true
        if (response.status === 204) {
          this.moreButton = false;
          this.Infomsg = "No more photos are available"
        }
      } catch (error) {
        this.handleError(error)
      }

    },
  },
  mounted() {
    this.loadContent();
  }
}
</script>

<template>


  <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
    <h1 class="h2">My stream</h1>

    <ErrorMsg v-if="errormsg" :msg="errormsg" class="mt-3"></ErrorMsg>
    <SuccessMsg v-if="successmsg" :msg="successmsg"></SuccessMsg>
  </div>
  <div class="justify-content-center align-items-center">

    <!-- Images -->
    <div class="d-grid justify-content-center" id="imageList">
      <Card v-bind:imageData="image" v-if="ready" v-for="image in photos"
            v-bind:username="this.username"></Card>
    </div>
    <div class="d-grid justify-content-center">
      <button @click="loadMorePhotos(this.page)" v-if="moreButton && ready" type="button" class="btn btn-outline-dark m-3">Load
        More
      </button>
      <InfoMsg :msg="this.Infomsg" v-if="this.Infomsg"></InfoMsg>
    </div>


  </div>

</template>

<style>
</style>
