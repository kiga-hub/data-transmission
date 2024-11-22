<template>
  <div class="layout">
    <Layout>
    <Header style="height:56px">
    </Header>
    <Content style="padding: 24px 50px">
      <Card>
        <p slot="title">
          <Icon type="ios-list"></Icon>
          Upgrade log of resources
        </p>
        <div slot="extra">
          <Button icon="md-sync" type="primary" @click="methodOnSynchronization">Synchronize data resources</Button>
          <Button icon="md-add" type="primary" @click="methodOnUpgrade">Upgrade data resources</Button>
        </div>
        <Table :columns="columns" :data="datas" max-height="1000" border>
          <template slot-scope="{ row }" slot="action">
            <Button type="success" size="small" @click="methodGetUpgradeLogDetail(row)">View</Button>
            <Button type="error" size="small" @click="methodDeleteUpgradeLogDetail(row)">Delete</Button>
          </template>
        </Table>
      </Card>
    </Content>
    <Footer class="layout-footer-center"></Footer>
  </Layout>
  <Modal
    v-model="synchronizationModal"
    title="Resource Synchronization"
    width="1024"
    @on-ok="methodOnClickSynchronization"
    >
    <Row type="flex" justify="center" :gutter="32">
      <Col style="margin: 20px;">
        <Form :model="synchronization_vars" label-position="left" ref="basic" :rules="synchronizationRuleValidate">
          <FormItem label="Server Synchronization Address" prop="url">
            <Input v-model="synchronization_vars.url" placeholder=""></Input>
          </FormItem>
        </Form>
      </Col>
    </Row>
  </Modal>
  <Modal
    v-model="upgradeModal"
    title="Resource Upgrade"
    width="1024"
    @on-ok="methodOnClickUpgrade"
    >
    <Row type="flex" justify="center" :gutter="32">
      <Col style="margin: 20px;">
        <Form :model="upgrade_vars" label-position="left" ref="basic" :rules="upgradeRuleValidate">
          <FormItem label="Server upgrade address" prop="remote_ip">
            <Input v-model="upgrade_vars.remote_ip" placeholder=""></Input>
          </FormItem>
          <FormItem label="Username" prop="user">
            <Input v-model="upgrade_vars.user" placeholder=""></Input>
          </FormItem>
          <FormItem label="Password" prop="password">
            <Input v-model="upgrade_vars.password" placeholder=""></Input>
          </FormItem>
          <FormItem label="Data Resource Name" prop="project_name">
            <Select v-model="source_list_index" @on-change="methodOnSourceListChange" :disabled="load_ok" placeholder="Prioritize synchronizing data resources">
              <Option v-for="(item, index) in sourceList" :key="item.project_name" :value="index">
                {{ item.project_name }}
              </Option>
            </Select>
          </FormItem>
        </Form>
      </Col>
    </Row>
  </Modal>

  <Modal
    v-model="detailModal"
    :title="detailTitle"
    width="1024"
    footer-hide>
    <Tabs value="log">
      <TabPane label="Logs" name="log">
        <Card v-if="detail.Log">
          <div class="markdown-container" id="log_div">
            <markdown-it-vue class="markdown-body" :content="detail.Log" :options="options" />
          </div>
        </Card>
      </TabPane>
      <TabPane label="Configure" name="config">
        <Card v-if="detail.Config">
          <div class="markdown-container">
            <json-viewer :value="detail.Config" class="markdown-body"></json-viewer>
          </div>
        </Card>
      </TabPane>
    </Tabs>
  </Modal>
</div>
</template>

<script>
import MarkdownItVue from 'markdown-it-vue'
import 'markdown-it-vue/dist/markdown-it-vue.css'
import {
  getUpgradeList, getUpgradeDetail, startUpgrade, deleteUpgradeDetail, updateSourceList, getSourceList
} from '@/api/data'

export default {
  name: 'Upgrade',
  components: {
    MarkdownItVue
  },
  data () {
    return {
      load_ok: true,
      source_list_index: -1,
      upgradeModal: false,
      synchronizationModal: false,
      upgrade_vars: {
        remote_ip: '192.168.8.244:22',
        user: 'root',
        password: 'password',
        project_name: ''
      },
      sourceList: [],
      synchronization_vars: {
        url: 'http://192.168.8.244:4567/download/source.json'
      },
      synchronizationRuleValidate: {
        url: [
          { required: true, message: 'The server synchronization address cannot be empty', trigger: 'blur' }
        ]
      },
      upgradeRuleValidate: {
        remote_ip: [
          { required: true, message: 'The server synchronization address cannot be empty', trigger: 'blur' }
        ],
        user: [
          { required: false, message: 'The username cannot be empty', trigger: 'blur' }
        ],
        password: [
          { required: false, message: 'The password cannot be empty', trigger: 'blur' }
        ],
        project_name: [
          { required: true, message: 'The data resource name cannot be empty', trigger: 'blur' }
        ]
      },
      columns: [{
        title: 'Data Resource Name',
        key: 'ProjectName'
      },
      {
        title: 'Upgrade Date',
        key: 'Date'
      },
      {
        title: 'Actions',
        slot: 'action',
        width: 150,
        align: 'center'
      }],
      datas: [],
      detail: {},
      detailModal: false,
      detailTitle: '',
      options: {
        markdownIt: {
          linkify: true,
          breaks: true
        },
        linkAttributes: {
          attrs: {
            target: '_blank',
            rel: 'noopener'
          }
        }
      },
      host: '192.168.8.244:22',
      isDevelopment: process.env.NODE_ENV !== 'production'
    }
  },
  created () {
    if (this.isDevelopment) {
      this.host = '192.168.8.244:22'
    } else {
      this.host = location.host
    }
  },
  mounted () {
    this.methodGetUpgradeLogList()
    this.methodFetchSourcelList()
  },
  methods: {
    methodOnUpgrade () {
      this.upgradeModal = true
      this.source_list_index = -1
      console.log('methodOnUpgrade', this.load_ok)
    },
    async methodFetchSourcelList () {
      try {
        const response = await getSourceList()
        this.sourceList = response.data.data
        console.log('Source list:', this.sourceList)
      } catch (error) {
        console.error('Failed to fetch source list:', error)
      }
    },
    methodOnSourceListChange (index) {
      console.log('methodOnSourceListChange', this.load_ok)
      if (this.sourceList.length === 0) {
        this.methodOnClickSynchronization().then(() => {
          this.updateProjectName(index)
        })
      } else {
        this.updateProjectName(index)
      }
    },
    updateProjectName (index) {
      if (index >= 0 && index < this.sourceList.length) {
        this.upgrade_vars.project_name = this.sourceList[index].project_name
        console.log('Project name updated to:', this.upgrade_vars.project_name)
      } else {
        console.error('Invalid index:', index)
      }
    },
    methodOnClickSynchronization () {
      this.methodUpdateSourcelList()
      console.log('methodOnClickSynchronization', this.load_ok)
    },
    methodUpdateSourcelList () {
      console.log('methodUpdateSourcelList called', this.synchronization_vars)
      updateSourceList(this.synchronization_vars).then(res => {
        if (res.data.code !== 200) {
          this.$Message.error(res.data.msg)
        }
        console.log('Source list updated:', res.data.data)
        this.sourceList = res.data.data
        this.load_ok = false
      })
      // window.location.reload()
    },
    methodStartUpgrade () {
      startUpgrade(this.upgrade_vars).then(res => {
        if (res.data.code !== 200) {
          this.$Message.error(res.data.msg)
          return
        }
        this.load_ok = true
        this.upgrade_vars.project_name = this.synchronization_vars.project_name
        // this.$Message.info(res.data.msg)
        this.methodGetUpgradeLogList()
      })
    },
    methodOnSynchronization () {
      this.load_ok = false
      this.synchronizationModal = true
    },
    methodOnClickUpgrade () {
      console.log('onClickUpgrade method called')
      console.log('Upgrade data:', this.upgrade_vars)
      this.methodStartUpgrade(this.upgrade_vars.project_name)
      setTimeout(() => {
        console.log('refresh current page!')
      }, 1000)
      this.$router.go(0)
    },
    methodGetUpgradeLogList () {
      var _this = this
      getUpgradeList().then(res => {
        if (res.data.code !== 200) {
          _this.$Message.error(res.data.msg)
          return
        }
        console.log('Upgrade log list:', res.data.data)

        // _this.datas = res.data.data
        _this.datas = res.data.data.map(item => {
          item.Date = item.Date.replace(
            /^(\d{4})(\d{2})(\d{2})(\d{2})(\d{2})(\d{2})$/,
            '$1-$2-$3 $4:$5:$6'
          )
          return item
        })
      })
    },
    methodGetUpgradeLogDetail (datas) {
      this.detailModal = true
      if (!this.detailModal) {
        return
      }

      console.log('Get upgrade log detail:', datas)

      var _this = this
      let params = { project_name: datas.ProjectName, date: datas.Date.replace(/^(\d{4})-(\d{2})-(\d{2}) (\d{2}):(\d{2}):(\d{2})$/, '$1$2$3$4$5$6') }

      getUpgradeDetail(params).then(res => {
        if (res.data.code !== 200) {
          _this.$Message.error(res.data.msg)
          return
        }

        _this.detail = res.data.data
        _this.detailTitle = 'Upgrade Time - ' + params.date
      })

      // setInterval(function () {
      //   _this.getUpgradeDetail(date)
      // }, 1000)
    },
    methodDeleteUpgradeLogDetail (datas) {
      var _this = this
      let params = { project_name: datas.ProjectName, date: datas.Date.replace(/^(\d{4})-(\d{2})-(\d{2}) (\d{2}):(\d{2}):(\d{2})$/, '$1$2$3$4$5$6') }

      deleteUpgradeDetail(params).then(res => {
        if (res.data.code !== 200) {
          _this.$Message.error(res.data.msg)
          return
        }
        _this.datas = _this.datas.filter(item => {
          return item.Date !== params.Date
        })
        _this.$Message.info(res.data.msg)
        _this.$router.go(0)
      })
    }
  }
}
</script>

<style scoped>
.layout{
    border: 1px solid #d7dde4;
    background: #f5f7f9;
    position: relative;
    border-radius: 4px;
    overflow: hidden;
}
.layout-footer-center{
    text-align: center;
}
.markdown-container {
  max-height: 500px; /* 设置最大高度，确保内容超出时显示滚动条 */
  overflow-y: auto; /* 设置为自动显示滚动条 */
}

.markdown-body {
  padding: 15px; /* 可选：为Markdown内容添加内边距 */
}
</style>
