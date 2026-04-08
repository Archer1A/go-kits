const fs = require('fs-extra');

const path = require('path');

const axios = require('axios');

// const ReadEnvFromRemote = require('./dist/src/read-dataSource').default;

// Interface representation (optional for documentation purposes)
// 调用异步函数并处理返回值或异常
// (async () => {
//     const url = 'http://localhost:8080/ind/cube/internal';  // 替换为实际的API URL
//     const modelName = '238';
//     const userId = '1783425296221384706';
//
//     try {
//         const data = await readFilters(url, modelName, userId);
//         console.log('Fetched data:', data);
//     } catch (error) {
//         console.error('Failed to fetch data:', error.message);
//     }
// })();

class IimpRepsitory {
    constructor(url, modelname) {
        this.url = url;
        this.modelname = modelname;
    }

    async dataSchemaFiles() {
        try {
            if (!this.modelname) {
                console.log("---modelname is empty,return....");
                return null;
            }
            console.log('---dataSchemaFiles--- ', `${this.url}/models/${this.modelname}`);
            const response = await axios.get(`${this.url}/models/${this.modelname}`);
            const files = response.data;
            if (process.env.CUBEJS_DEV_MODE && process.env.CUBEJS_DEV_MODE === 'true') {
                const schemaPath = process.env.CUBEJS_SCHEMA_PATH || 'model';
                console.log('schemaPath is ', schemaPath);
                if (!fs.existsSync(schemaPath)) {
                    console.log(`schemaPath : ${schemaPath},is not exist`);
                    fs.mkdirSync(schemaPath);
                    fs.mkdirSync(path.join(schemaPath, 'cubes'));
                    fs.mkdirSync(path.join(schemaPath, 'views'));
                }
                await fs.emptyDir(path.join(schemaPath, 'cubes'));
                await fs.emptyDir(path.join(schemaPath, 'views'));
                await Promise.all(files.map(file => fs.writeFile(path.join(schemaPath, 'cubes', file.fileName), file.content)));
            }
            return files;
        } catch (error) {
            console.error('Error fetching data:', error.message);
            throw error;
        }
    }

    // protected saveToModel (implementation would be similar)
}

async function listModels(url){
    try {
        console.log('read query models from ', `${url}/models`);
        const response = await axios.get(`${url}/models`);
        return response.data;
    } catch (error) {
        console.error('Error fetching data:', error.message);
        throw error;
    }
}
async function readFromUrl(url, dataSource) {
    try {
        console.log('read data source from ', `${url}/dataSource/${dataSource}`);
        const response = await axios.get(`${url}/dataSource/${dataSource}`);
        return response.data;
    } catch (error) {
        console.error('Error fetching data:', error.message);
        throw error;
    }
}

async function readSchemaVersionfromUrl(url, name) {
    try {
        if (!name) {
            return null;
        }
        console.log('read schema version from ', `${url}/models/${name}/version`);
        const response = await axios.get(`${url}/models/${name}/version`);
        return response.data;
    } catch (error) {
        console.error('Error fetching data:', error.message);
        throw error;
    }
}

async function readFilters(url,modelName,userId){
    try {
        if (!modelName) {
            console.log("readFilters modelname is null")
            return null;
        }
        console.log('read query overwrite from ', `${url}/models/${modelName}/filters/${userId}`);
        const response = await axios.get(`${url}/models/${modelName}/filters/${userId}`);
        return response.data;
    } catch (error) {
        console.error('Error fetching data:', error.message);
        throw error;
    }
}

console.log(IimpRepsitory);
module.exports = {
    repositoryFactory: ({ securityContext }) => {
        const url = process.env.REMOTE_MODLE_URL || 'http://host.docker.internal:8080/ind/cube/internal';
        const modelname = securityContext.modelname;
        console.log(`load from ${securityContext},load from remote from ${url}`);

        return new IimpRepsitory(url, modelname);
    },
    driverFactory: ({ securityContext, dataSource }) => {
        const url = process.env.REMOTE_DATASOURCE_URL || 'http://host.docker.internal:8080/ind/cube/internal';
        const ds = dataSource || 'default';
        console.log(`load EVN from  ${url}, dataSource is ${ds}`);
        return readFromUrl(url, ds);
    },

    queryRewrite: async (query, { securityContext }) => {
        const model_name = securityContext.modelname
        const user_id = securityContext.user_id
        const url =  process.env.REMOTE_MODLE_URL || 'http://host.docker.internal:8080/ind/cube/internal';
        const filters = await readFilters(url,model_name,user_id);
        for (const filter of filters){
            query.filters.push(filter)
        }
        return query;
    },

    schemaVersion: ({ securityContext }) => {

        const url = process.env.REMOTE_MODELVERSION_URL || 'http://host.docker.internal:8080/ind/cube/internal';
        const modelname = securityContext.modelname;
        if (url) {
            return readSchemaVersionfromUrl(url, modelname);
        } else {
            return 1;
        }
    },

    // extendContext: ({ securityContext }) => {
    //     console.log("0000.............")
    //     console.log(securityContext)
    //     console.log("2222.............")
    //
    //     return {
    //         securityContext: {
    //             modelname: '1783425296221384706_273',
    //         }
    //     };
    // },

    checkSqlAuth: async  (req, username,password) => {
        const splitStrings = username.split('_');
        const returnObj = {
            password,
            securityContext: {
                modelname: splitStrings[1],
                user_id: splitStrings[0],
            },
        };
        return returnObj;
    },



    contextToAppId: ({ securityContext }) => {
        console.log("contextToAppId")
        console.log(securityContext)
        // `CUBE_APP_${securityContext.tenant_id}_${securityContext.modelname}`,
        return `CUBE_APP_${securityContext.modelname}`;
    },

    scheduledRefreshContexts: async () =>{
        const url =  process.env.REMOTE_MODLE_URL || 'http://host.docker.internal:8080/ind/cube/internal';
        let newVar = await listModels(url);
        let contextList = [];
        for(const id of newVar){
            let  a = {
                securityContext: {
                    modelname: String(id)
                }
            }
            contextList.push(a);
        }
        return contextList;
    } ,

    contextToApiScopes: (securityContext, defaultScopes) => ['meta', 'data', 'graphql', 'jobs']
};
