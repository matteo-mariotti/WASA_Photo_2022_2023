<script>
import backend from "@/services/axios";
import ErrorMsg from "@/components/ErrorMsg.vue";
import SuccessMsg from "@/components/SuccessMsg.vue";
import UserSummaryCard from "@/components/UserSummaryCard.vue";
import InfoMsg from "@/components/InfoMsg.vue";

export default {
  components: {InfoMsg, UserSummaryCard, SuccessMsg, ErrorMsg},
  data: function () {
    return {
      errormsg: null,
      loading: false,
      identifier: null,
      searchBox: null,
      userList:[],
      Infomsg: null,
    }
  },
  methods: {
    async search(){
      try{
      if (this.searchBox.length >=3) {
        this.Infomsg = null
        this.userList = []
        let response = await backend.get(`/users?username=${this.searchBox}`)
        if (response.status === 200) {
          response.data.forEach(user => {
            this.userList.push(user)
          })
        }
      }else{
        this.userList = []
        this.Infomsg = "Write at least three characters"
      }
      }catch (e) {
        this.errormsg = e.response.message
        console.log(e)
      }
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

        <input v-on:input="search" type="text" v-model="searchBox" class="form-control mb-5" placeholder="username">


        <div>
          <InfoMsg :msg="this.Infomsg" v-if="this.Infomsg"></InfoMsg>
          <UserSummaryCard v-for="user in this.userList" :username="user"></UserSummaryCard>

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
