<script>
import backend from "@/services/axios";
import ErrorMsg from "@/components/ErrorMsg.vue";
import UserSummaryCard from "@/components/UserSummaryCard.vue";
import InfoMsg from "@/components/InfoMsg.vue";

export default {
  components: {InfoMsg, UserSummaryCard, ErrorMsg},
  data: function () {
    return {
      errormsg: null,
      loading: false,
      identifier: null,
      searchBox: null,
      userList:[],
      Infomsg: null,
      page: 0,
      moreButton: false,
    }
  },
  methods: {
    async search(){
      try{
      if (this.searchBox.length >=3) {
        this.Infomsg = null
        let response = await backend.get(`/users?username=${this.searchBox}&page=${this.page}`)
        this.page = this.page +1
        if (response.status === 200) {
          response.data.forEach(user => {
            this.userList.push(user)
          })
        }else if (response.status === 204){
          this.Infomsg = "No users found"
          if (this.moreButton){
            this.moreButton = false
          }
        }

      }else{
        this.moreButton = false
        this.userList = []
        this.Infomsg = "Write at least three characters"
      }
      }catch (e) {
        this.errormsg = e.response.message
        console.log(e)
      }
    },
    async searchNew(){
      this.moreButton = true
      this.page = 0
      this.userList = []
      await this.search()
    }

  },
  mounted() {
    this.Infomsg = "Write at least three characters"
  }
}

</script>

<template>


  <div class="container mt-5">
    <div class="row">
      <div class="col-sm-3"></div>
      <div class="col-sm-6">



        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

        <div style="text-align: center">
          <h5>Search a user</h5>
          <hr>
        </div>

        <input v-on:input="searchNew" type="text" v-model="searchBox" class="form-control mb-5" placeholder="username">


        <div>
          <UserSummaryCard v-for="user in this.userList" v-bind:key="user" :username="user"></UserSummaryCard>

          <InfoMsg :msg="this.Infomsg" v-if="this.Infomsg"></InfoMsg>
          <div v-if="moreButton" style="text-align: center">
          <button @click="search" type="button" class="btn btn-outline-dark m-3">Load
            More
          </button>
          </div>
        </div>

      </div>

      <div class="col-sm-3"></div>

    </div>

  </div>
</template>

<style>
::placeholder {
  text-align: center;
}
</style>
