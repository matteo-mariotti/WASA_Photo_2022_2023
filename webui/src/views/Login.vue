<script>
import backend from "@/services/axios";
import ErrorMsg from "@/components/ErrorMsg.vue";
import SuccessMsg from "@/components/SuccessMsg.vue";

export default {
  components: {SuccessMsg, ErrorMsg},
  data: function () {
    return {
      errormsg: null,
      loading: false,
      identifier: null,
      successmsg: null,
    }
  },
  methods: {
    async sendData() {
      try {
        let response = await backend.post("/session", {
          identifier: this.identifier,
        });
        console.log(response)
        this.handleResponse(response);
      } catch (error) {
        console.log(error)
        this.handleError(error);
      }
    },
    handleResponse(response) {
      // Check if the response is 200 (OK)
      if (response.status === 200) {
        // If the response is 200, then the user is logged in
        // We can redirect him to the home page
        sessionStorage.setItem("token", response.data.token);
        sessionStorage.setItem("username", this.identifier);
        sessionStorage.setItem("logged", true);
        this.errormsg = null
        this.successmsg = "Login successful"
        this.loading = true
        this.$emit("logging")
        setTimeout(this.redirect, 1500)
      }
    },
    redirect(){
      this.$router.push("/");
    },
    handleError(error) {
      // Print the error from the server (for debugging)
      this.errormsg = error.response.data.message
    },
  },
}
</script>

<template>


  <div class="container mt-5">
    <div class="row">
      <div class="col-sm-3"></div>
      <div class="col-sm-6">


        <div style="text-align: center">
          <h2>Welcome to WASA Photo</h2>
          <br>
        </div>

         <input type="text" v-model="identifier" class="form-control" placeholder="username">
          <div style="text-align: center">
            <button @click="sendData()" type="submit" class="btn btn-primary m-3">Sign in/up</button>
          </div>
        <br>

        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        <SuccessMsg v-if="successmsg" :msg="successmsg"></SuccessMsg>

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
