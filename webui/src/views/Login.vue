<script>
import backend from "@/services/axios";
import ErrorMsg from "@/components/ErrorMsg.vue";
import SuccessMsg from "@/components/SuccessMsg.vue";

export default {
  emits: ["logging"],
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
    valid(){
      if (this.identifier.length <3) {
        document.getElementById("logButton").disabled = true;
        this.errormsg = "Usernames must have at least three characters"
        return false
      }
      if (this.identifier.length >12) {
        document.getElementById("logButton").disabled = true;
        this.errormsg = "Usernames can't have more than 12 characters"
        return false
      }
      this.errormsg = null
      document.getElementById("logButton").disabled = false;
      return true
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
        <form @submit.prevent="sendData()">
         <input type="text" v-model="identifier" v-on:input="valid" class="form-control" placeholder="username">
          <div style="text-align: center">
            <button id="logButton" type="submit" class="btn btn-primary m-3" disabled>Sign in/up</button>
          </div>
          </form>
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
