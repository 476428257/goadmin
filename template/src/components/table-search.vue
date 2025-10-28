<template>
	<div class="search-container" v-if="showsearch">
		<el-form ref="searchRef" :model="query" :inline="true" >
			<el-form-item :label="item.label" :prop="item.prop" v-for="item in options">
				<!-- 文本框、下拉框、日期框 -->
				<el-input v-if="item.type === 'input'" v-model="query[item.prop]" :disabled="item.disabled"
					:placeholder="item.placeholder" clearable style="width: 150px"></el-input>
				<el-select v-else-if="item.type === 'select'" v-model="query[item.prop]" :disabled="item.disabled"
					:placeholder="item.placeholder" style="width: 130px" clearable>
					<el-option v-for="opt in item.opts" :label="opt.label" :value="opt.value"></el-option>
				</el-select>
				<el-date-picker v-else-if="item.type === 'date'" type="datetime" v-model="query[item.prop]"
					:value-format="item.format"></el-date-picker>
				<el-date-picker v-else-if="item.type === 'daterange'" type="datetimerange" v-model="query[item.prop]"
						:value-format="item.format" range-separator="至" start-placeholder="开始日期" 
						end-placeholder="结束日期" :disabled="item.disabled" clearable></el-date-picker>	
			</el-form-item>
			<el-form-item>
				<el-button type="primary" :icon="Search" @click="search">搜索</el-button>
				<el-button :icon="Refresh" @click="resetForm(searchRef)">重置</el-button>
			</el-form-item>
		</el-form>
	</div>
</template>

<script lang="ts" setup>
import { FormOptionList } from '@/types/form-option';
import { Refresh, Search } from '@element-plus/icons-vue';
import { FormInstance } from 'element-plus';
import { PropType, ref } from 'vue';

const props = defineProps({
	showsearch: {
		type: Boolean,
		default: false
	},
	query: {
		type: Object,
		required: true
	},
	options: {
		type: Array as PropType<Array<FormOptionList>>,
		required: true
	},
	search: {
		type: Function,
		default: () => { }
	}
});

const searchRef = ref<FormInstance>();
const resetForm = (formEl: FormInstance | undefined) => {
	if (!formEl) return
	formEl.resetFields()
	props.search();
}
</script>

<style scoped>
.search-container {
	padding: 20px 30px 0;
	background-color: #fff;
	margin-bottom: 10px;
	border: 1px solid #ddd;
	border-radius: 5px
}
</style>
