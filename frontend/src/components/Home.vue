<template>
  <v-card>
    <v-toolbar prominent color="primary">
      <v-app-bar-nav-icon></v-app-bar-nav-icon>

      <v-toolbar-title>MamuroEmail</v-toolbar-title>

      <v-spacer></v-spacer>
    </v-toolbar>
  </v-card>
  <v-container fluid>
    <v-card-title class="mx-auto" densed>
      <v-spacer></v-spacer>
      <v-text-field v-model="search" v-on:change="searchEmail" append-icon="mdi-magnify" label="Search"
        density="compact" single-line hide-details>
      </v-text-field>
    </v-card-title>
    <v-row no-gutters>
      <v-col cols="12" sm="7" maxheight="460">
        <v-card>
          <v-card-table width="500">
            <v-table height="460">
              <thead>
                <tr>
                  <th class="text-left">
                    Subject
                  </th>
                  <th class="text-left">
                    From
                  </th>
                  <th class="text-left">
                    To
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="message in messages" :key="message" v-on:click="showEmail(message._source.Message)">
                  <td>{{message._source.Subject}}</td>
                  <td>{{message._source.From}}</td>
                  <td>{{message._source.To}}</td>
                </tr>
              </tbody>
            </v-table>
          </v-card-table>
        </v-card>
      </v-col>
      <v-col cols="12" sm="5">
        <v-textarea 
          variant="solo"
          :value="cMessage"
          height="460"
          rows="18"
          no-resize
          readonly
        ></v-textarea>
      </v-col>
    </v-row>
  </v-container>

</template>

<script>
export default {
  /* eslint-disable */
  name: 'Home',
  data() {
    return {
      messages: [],
      cMessage: '',
      search: ''
    }
  },
  methods: {
    async searchEmail() {
      let options = {
        method: 'POST',
        mode: 'cors',
        headers: {
          'Access-Control-Allow-Origin': '*',
          'Access-Control-Request-Method': 'POST, OPTIONS'
        },
        body: JSON.stringify({ input: this.search })
      }
      this.messages=[];
      const resp = fetch("http://localhost:3000/searchMail", options);
      resp.then(res => res.json()).then(data => data.hits.hits).then(hits => hits.map((data) => this.messages.push(data)))
    },
    showEmail(data) {
      this.cMessage = data
      this.$refs.update.blur()
      console.log(this.$refs.update)
    }
  }
}
</script>