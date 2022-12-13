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
        this.handleResponse(response);
      } catch (error) {
        this.handleError(error);
      }
    },
    handleResponse(response) {
      // Check if the response is 200 (OK)
      if (response.status === 200) {
        // If the response is 200, then the user is logged in
        // We can redirect him to the home page
        sessionStorage.setItem("token", response.data.token);
        this.errormsg = null
        this.successmsg = "Login successful"
        this.loading = true
        setTimeout(this.redirect, 1500)
      }
    },
    redirect(){
      this.$router.push("/");
    },
    handleError(error) {

      // Print the error from the server (for debugging)
      this.errormsg = error.response.data.message

      // Log the error on console
      console.log(error)

    },
  },
}
</script>

<template>


  <div class="container mt-5">
    <div class="row">
      <div class="col-3"></div>
      <div class="col-6">


        <div style="text-align: center">
          <label>Inserisci il tuo username per accedere</label>
          <br>
        </div>

        <form>
          <div class="form-group">
            <input type="text" v-model="identifier" class="form-control" placeholder="username">
          </div>
          <br>
          <div style="text-align: center">
            <button @click="sendData()" type="submit" class="btn btn-primary">Sign in</button>
          </div>
        </form>

        <br>

        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        <SuccessMsg v-if="successmsg" :msg="successmsg"></SuccessMsg>

      </div>

    </div>

  </div>
</template>

<style>
</style>
