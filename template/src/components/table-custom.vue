<template>
    <div>
        <div class="table-toolbar" v-if="hasToolbar">
            <div class="table-toolbar-left btn-group">
                <el-button type="info" style="background-color: #30364a;" :icon="Refresh"  @click="refresh"></el-button>
                <el-button type="success" v-if="hasadd" :icon="CirclePlusFilled" @click="handleAdd">新增</el-button>
                <el-button type="warning" v-if="hasview" :icon="View" :disabled="multipleSelection.length === 0" @click="viewFunc(multipleSelection)">查看</el-button>
                <el-button type="primary" v-if="hasedit" :icon="Edit" :disabled="multipleSelection.length === 0" @click="editFunc(multipleSelection)">编辑</el-button>
                <el-button type="danger" v-if="hasdel" :icon="Delete" :disabled="multipleSelection.length === 0" @click="handleDelSelection(multipleSelection)">删除</el-button>
                <slot name="toolbarBtn"></slot>
            </div>
            <div class="table-toolbar-right flex-center">
                <el-tooltip effect="dark" content="搜索" placement="top" v-if="showsearchbtn">
                    <el-icon class="columns-setting-icon" @click="toggleSearch">
                        <Search />
                    </el-icon>
                </el-tooltip>
                <el-divider direction="vertical" v-if="m_d1" />
                <el-tooltip effect="dark" content="列设置" placement="top" v-if="setbtn">
                    <el-dropdown :hide-on-click="false" size="small" trigger="click">
                        <el-icon class="columns-setting-icon">
                            <Setting />
                        </el-icon>
                        <template #dropdown>
                            <el-dropdown-menu>
                                <el-dropdown-item v-for="c in columns">
                                    <el-checkbox v-model="c.visible" :label="c.label" />
                                </el-dropdown-item>
                            </el-dropdown-menu>
                        </template>
                    </el-dropdown>
                </el-tooltip>
            </div>
        </div>
        <el-table class="mgb20" :style="{ width: '100%' }" border :data="tableData" :row-key="rowKey"
            @selection-change="handleSelectionChange" ref="tableRef"
            @row-click="handleRowClickEvent"  :row-class-name="rowClassName" table-layout="auto"
            :default-expand-all="defaultExpandAll" :tree-props="treeProps" >
            <template v-for="item in columns" :key="item.prop">
                <el-table-column v-if="item.visible" :prop="item.prop" :label="item.label" :width="item.width"
                    :type="item.type" :align="item.align || 'center'">

                    <template #default="{ row, column, $index }" v-if="item.type === 'index'">
                        {{ getIndex($index) }}
                    </template>
                    <template #default="{ row, column, $index }" v-if="!item.type">
                        <slot :name="item.prop" :rows="row" :index="$index">
                            <template v-if="item.prop == 'operator'">
                                
                            </template>
                            <template v-else-if="item.isImage">
                                <div class="image-cell">
                                    <template v-if="toImageArray(row[item.prop]).length">
                                        <img
                                            v-for="(src, idx) in toImageArray(row[item.prop])"
                                            :key="idx"
                                            :src="src"
                                            class="table-image"
                                            :alt="item.label"
                                            @click.stop="showImagePreview(toImageArray(row[item.prop]), idx)"
                                        />
                                    </template>
                                    <span v-else class="no-image">暂无图片</span>
                                </div>
                            </template>
                            <template v-else-if="item.ellipsisClick">
                                <div 
                                    class="ellipsis-cell" 
                                    @click="showContentDetail(row[item.prop])"
                                    :title="row[item.prop]"
                                >
                                    {{ row[item.prop] }}
                                </div>
                            </template>
                            <span v-else-if="item.formatter">
                                {{ item.formatter(row[item.prop]) }}
                            </span>
                            <span v-else>
                                {{ row[item.prop] }}
                            </span>
                        </slot>
                    </template>
                </el-table-column>
            </template>
        </el-table>
        <el-pagination v-if="hasPagination" :current-page="currentPage" :page-size="pageSize" :background="true"
            :layout="layout" :total="total" @current-change="handleCurrentChange" />
        
        <!-- 内容详情查看弹窗 -->
        <el-dialog title="内容详情" v-model="contentDetailVisible" width="800px" destroy-on-close>
            <div class="content-detail">
                <el-input
                    v-model="contentDetailData"
                    type="textarea"
                    :rows="15"
                    readonly
                    resize="none"
                />
            </div>
        </el-dialog>
        
        <!-- 图片预览弹窗 -->
        <el-dialog title="图片预览" v-model="imagePreviewVisible" width="60%" destroy-on-close>
            <div class="image-preview">
                <div class="preview-controls" v-if="previewImages.length">
                    <el-button size="small" @click="prevImage" :disabled="previewIndex <= 0">上一张</el-button>
                    <span class="preview-counter">{{ previewIndex + 1 }} / {{ previewImages.length }}</span>
                    <el-button size="small" @click="nextImage" :disabled="previewIndex >= previewImages.length - 1">下一张</el-button>
                </div>
                <img 
                    v-if="currentPreviewImage" 
                    :src="currentPreviewImage" 
                    class="preview-image" 
                    :alt="'图片预览'"
                />
                <div v-else class="no-image-preview">暂无图片</div>
            </div>
        </el-dialog>
    </div>
</template>

<script setup lang="ts">

import { CirclePlusFilled, Delete, Edit, Refresh, Search, View } from '@element-plus/icons-vue';
import { ElMessageBox } from 'element-plus';
import { computed, PropType, ref, toRefs } from 'vue';
const tableRef = ref();
const handleRowClickEvent = (row, column, event) => {
  tableRef.value.toggleRowSelection(row);
};
const rowClassName = ({ row }) => {
  // 如果当前行在多选数组中，则加高亮类
  return multipleSelection.value.find(item => item.id === row.id) ? 'multi-row-highlight' : '';
};
const showsearch = defineModel('showsearch');
const toggleSearch = () => {
    showsearch.value = !showsearch.value;
};
const props = defineProps({
    showsearchbtn: {
        type: Boolean,
        default: true
    },
    m_d1: {
        type: Boolean,
        default: true
    },
    setbtn: {
        type: Boolean,
        default: true
    },
    hasadd: {
        type: Boolean,
        default: true
    },
    hasview: {
        type: Boolean,
        default: true
    },
    hasedit: {
        type: Boolean,
        default: true
    },
    hasdel: {
        type: Boolean,
        default: true
    },
    showsearch: {
        type: Boolean,
        default: false
    },
    // 表格相关
    tableData: {
        type: Array,
        default: []
    },
    columns: {
        type: Array as PropType<any[]>,
        default: []
    },
    rowKey: {
        type: String,
        default: 'id'
    },
    hasToolbar: {
        type: Boolean,
        default: true
    },
    //  分页相关
    hasPagination: {
        type: Boolean,
        default: true
    },
    total: {
        type: Number,
        default: 0
    },
    currentPage: {
        type: Number,
        default: 1
    },
    pageSize: {
        type: Number,
        default: 10
    },

    layout: {
        type: String,
        default: 'total, prev, pager, next'
    },
    viewFunc: {
        type: Function,
        default: () => { }
    },
    editFunc: {
        type: Function,
        default: () => { }
    },
    delSelection: {
        type: Function,
        default: () => { }
    },
    refresh: {
        type: Function,
        default: () => { }
    },
    handleAdd: {
        type: Function,
        default: () => { }
    },
    changePage: {
        type: Function,
        default: () => { }
    },
    // 树形表格相关
    defaultExpandAll: {
        type: Boolean,
        default: false
    },
    treeProps: {
        type: Object,
        default: () => ({ children: 'children', hasChildren: 'hasChildren' })
    },
    checkSelectable: {
        type: Function,
        default: null
    }
})

let {
    tableData,
    columns,
    rowKey,
    hasToolbar,
    hasPagination,
    total,
    currentPage,
    pageSize,
    layout,
    defaultExpandAll,
    treeProps,
    checkSelectable,
} = toRefs(props)

columns.value.forEach((item) => {
    if (item.visible === undefined) {
        item.visible = true
    }
})

// 当选择项发生变化时会触发该事件
const multipleSelection = ref([])
const handleSelectionChange = (selection: any[]) => {
    multipleSelection.value = selection
}

// 内容详情弹窗相关
const contentDetailVisible = ref(false);
const contentDetailData = ref('');

// 图片预览弹窗相关（支持数组与切换）
const imagePreviewVisible = ref(false);
const previewImages = ref<string[]>([]);
const previewIndex = ref(0);
const currentPreviewImage = computed(() => previewImages.value[previewIndex.value] || '');

// 当前页码变化的事件
const handleCurrentChange = (val: number) => {
    props.changePage(val)
}

const handleDelSelection = (row) => {
    ElMessageBox.confirm('确定要删除吗？', '提示', {
        type: 'warning'
    })
        .then(async () => {
            props.delSelection(row);
        })
        .catch(() => { });
};

const getIndex = (index: number) => {
    return index + 1 + (currentPage.value - 1) * pageSize.value
}

// 显示内容详情
const showContentDetail = (content: string) => {
    contentDetailData.value = content || '无数据';
    contentDetailVisible.value = true;
};

// 将逗号分隔的图片字符串转为数组，已是数组则直接返回
const toImageArray = (val: any): string[] => {
    if (!val) return [];
    if (Array.isArray(val)) {
        return val.filter(Boolean);
    }
    if (typeof val === 'string') {
        return val
            .split(',')
            .map(s => s.trim())
            .filter(s => s.length > 0);
    }
    return [];
};

// 显示图片预览（数组 + 指定起始下标）
const showImagePreview = (images: string[] | string, startIndex = 0) => {
    const arr = Array.isArray(images) ? images : toImageArray(images);
    previewImages.value = arr;
    previewIndex.value = Math.min(Math.max(startIndex, 0), Math.max(arr.length - 1, 0));
    imagePreviewVisible.value = true;
};

const prevImage = () => {
    if (previewIndex.value > 0) previewIndex.value -= 1;
};

const nextImage = () => {
    if (previewIndex.value < previewImages.value.length - 1) previewIndex.value += 1;
};

// 展开/收缩所有树形表格节点
const toggleExpandAll = (expand: boolean) => {
    if (!tableRef.value) return;
    
    const toggleRow = (data: any[]) => {
        data.forEach(row => {
            if (row.children && row.children.length > 0) {
                tableRef.value.toggleRowExpansion(row, expand);
                toggleRow(row.children);
            }
        });
    };
    
    toggleRow(tableData.value);
};

// 暴露方法给父组件
defineExpose({
    toggleExpandAll
});

</script>

<style scoped>
.table-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 10px;
}

.columns-setting-icon {
    display: block;
    font-size: 18px;
    cursor: pointer;
    color: #676767;
}

.ellipsis-cell {
    max-width: 300px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    cursor: pointer;
    color: #409eff;
    transition: color 0.3s;
}

.ellipsis-cell:hover {
    color: #66b1ff;
    text-decoration: underline;
}

.content-detail {
    max-height: 400px;
    overflow-y: auto;
}

.content-detail .el-textarea__inner {
    font-family: 'Courier New', monospace;
    font-size: 13px;
    line-height: 1.4;
}

/* 图片相关样式 */
.image-cell {
    cursor: pointer;
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 4px;
}

.table-image {
    width: 60px;
    height: 60px;
    object-fit: cover;
    border-radius: 4px;
    border: 1px solid #dcdfe6;
    transition: transform 0.3s ease;
}

.table-image:hover {
    transform: scale(1.1);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.no-image {
    color: #909399;
    font-size: 12px;
    padding: 20px 0;
}

.image-preview {
    text-align: center;
    padding: 20px;
}

.preview-image {
    max-width: 100%;
    max-height: 70vh;
    object-fit: contain;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.no-image-preview {
    color: #909399;
    font-size: 16px;
    padding: 60px 0;
}
</style>
<style>
.table-header .cell {
    color: #333;
}
.multi-row-highlight {
  background: #e6f7ff !important; 
}
.btn-group .el-button {
  margin-right: -8px !important; /* 默认是12px，改小一半 */
}

</style>