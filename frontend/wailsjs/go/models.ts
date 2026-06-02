export namespace main {
	
	export class BucketView {
	    name: string;
	    location: string;
	    creationDate: string;
	    storageClass: string;
	    extranetEndpoint: string;
	    intranetEndpoint: string;
	
	    static createFrom(source: any = {}) {
	        return new BucketView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.location = source["location"];
	        this.creationDate = source["creationDate"];
	        this.storageClass = source["storageClass"];
	        this.extranetEndpoint = source["extranetEndpoint"];
	        this.intranetEndpoint = source["intranetEndpoint"];
	    }
	}
	export class CMSMetricsView {
	    instanceId: string;
	    cpuUtilization?: number;
	    memoryUtilization?: number;
	    diskReadBps?: number;
	    diskWriteBps?: number;
	    internetRx?: number;
	    internetTx?: number;
	    updateTime: string;
	
	    static createFrom(source: any = {}) {
	        return new CMSMetricsView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.instanceId = source["instanceId"];
	        this.cpuUtilization = source["cpuUtilization"];
	        this.memoryUtilization = source["memoryUtilization"];
	        this.diskReadBps = source["diskReadBps"];
	        this.diskWriteBps = source["diskWriteBps"];
	        this.internetRx = source["internetRx"];
	        this.internetTx = source["internetTx"];
	        this.updateTime = source["updateTime"];
	    }
	}
	export class CMSMetricsResult {
	    success: boolean;
	    message: string;
	    metrics: CMSMetricsView[];
	
	    static createFrom(source: any = {}) {
	        return new CMSMetricsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.metrics = this.convertValues(source["metrics"], CMSMetricsView);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class InstanceView {
	    instanceId: string;
	    instanceName: string;
	    status: string;
	    regionId: string;
	    zoneId: string;
	    publicIp: string;
	    privateIp: string;
	    creationTime: string;
	
	    static createFrom(source: any = {}) {
	        return new InstanceView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.instanceId = source["instanceId"];
	        this.instanceName = source["instanceName"];
	        this.status = source["status"];
	        this.regionId = source["regionId"];
	        this.zoneId = source["zoneId"];
	        this.publicIp = source["publicIp"];
	        this.privateIp = source["privateIp"];
	        this.creationTime = source["creationTime"];
	    }
	}
	export class LogEntryView {
	    timestamp: number;
	    content: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new LogEntryView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timestamp = source["timestamp"];
	        this.content = source["content"];
	    }
	}
	export class RegionResultView {
	    region: string;
	    instances: InstanceView[];
	    totalCount: number;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new RegionResultView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.region = source["region"];
	        this.instances = this.convertValues(source["instances"], InstanceView);
	        this.totalCount = source["totalCount"];
	        this.error = source["error"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MultiRegionResult {
	    success: boolean;
	    message: string;
	    regions: RegionResultView[];
	
	    static createFrom(source: any = {}) {
	        return new MultiRegionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.regions = this.convertValues(source["regions"], RegionResultView);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class OSSBucketResult {
	    success: boolean;
	    message: string;
	    buckets: BucketView[];
	
	    static createFrom(source: any = {}) {
	        return new OSSBucketResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.buckets = this.convertValues(source["buckets"], BucketView);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ObjectView {
	    key: string;
	    size: number;
	    lastModified: string;
	    etag: string;
	    type: string;
	    storageClass: string;
	    isFolder: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ObjectView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.size = source["size"];
	        this.lastModified = source["lastModified"];
	        this.etag = source["etag"];
	        this.type = source["type"];
	        this.storageClass = source["storageClass"];
	        this.isFolder = source["isFolder"];
	    }
	}
	export class OSSObjectResult {
	    success: boolean;
	    message: string;
	    objects: ObjectView[];
	    isTruncated: boolean;
	
	    static createFrom(source: any = {}) {
	        return new OSSObjectResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.objects = this.convertValues(source["objects"], ObjectView);
	        this.isTruncated = source["isTruncated"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class OperationResult {
	    success: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new OperationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	    }
	}
	export class QueryECSResult {
	    success: boolean;
	    message: string;
	    instances: InstanceView[];
	    totalCount: number;
	
	    static createFrom(source: any = {}) {
	        return new QueryECSResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.instances = this.convertValues(source["instances"], InstanceView);
	        this.totalCount = source["totalCount"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class SLSLogQueryResult {
	    success: boolean;
	    message: string;
	    logs: LogEntryView[];
	    count: number;
	    hasMore: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SLSLogQueryResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.logs = this.convertValues(source["logs"], LogEntryView);
	        this.count = source["count"];
	        this.hasMore = source["hasMore"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SLSLogStoreResult {
	    success: boolean;
	    message: string;
	    logStores: string[];
	
	    static createFrom(source: any = {}) {
	        return new SLSLogStoreResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.logStores = source["logStores"];
	    }
	}
	export class SLSStreamResult {
	    success: boolean;
	    message: string;

	    static createFrom(source: any = {}) {
	        return new SLSStreamResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	    }
	}
	export class CloudProductView {
	    id: string;
	    name: string;
	    namespace: string;
	    metrics: MetricView[];

	    static createFrom(source: any = {}) {
	        return new CloudProductView(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.namespace = source["namespace"];
	        this.metrics = this.convertValues(source["metrics"], MetricView);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CloudProductResult {
	    success: boolean;
	    message: string;
	    products: CloudProductView[];

	    static createFrom(source: any = {}) {
	        return new CloudProductResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.products = this.convertValues(source["products"], CloudProductView);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MetricView {
	    id: string;
	    name: string;
	    unit: string;
	    description: string;

	    static createFrom(source: any = {}) {
	        return new MetricView(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.unit = source["unit"];
	        this.description = source["description"];
	    }
	}
	export class VPCView {
	    vpcId: string;
	    vpcName: string;
	    cidrBlock: string;
	    status: string;
	    regionId: string;
	    description: string;
	    creationTime: string;

	    static createFrom(source: any = {}) {
	        return new VPCView(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vpcId = source["vpcId"];
	        this.vpcName = source["vpcName"];
	        this.cidrBlock = source["cidrBlock"];
	        this.status = source["status"];
	        this.regionId = source["regionId"];
	        this.description = source["description"];
	        this.creationTime = source["creationTime"];
	    }
	}
	export class VPCDetailView {
	    vpcId: string;
	    vpcName: string;
	    cidrBlock: string;
	    status: string;
	    regionId: string;
	    description: string;
	    creationTime: string;
	    vswitchIds: string[];
	    natGatewayIds: string[];
	    routerTableIds: string[];

	    static createFrom(source: any = {}) {
	        return new VPCDetailView(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vpcId = source["vpcId"];
	        this.vpcName = source["vpcName"];
	        this.cidrBlock = source["cidrBlock"];
	        this.status = source["status"];
	        this.regionId = source["regionId"];
	        this.description = source["description"];
	        this.creationTime = source["creationTime"];
	        this.vswitchIds = source["vswitchIds"];
	        this.natGatewayIds = source["natGatewayIds"];
	        this.routerTableIds = source["routerTableIds"];
	    }
	}
	export class VPCDetailResult {
	    success: boolean;
	    message: string;
	    detail?: VPCDetailView;

	    static createFrom(source: any = {}) {
	        return new VPCDetailResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.detail = this.convertValues(source["detail"], VPCDetailView);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class VPCListResult {
	    success: boolean;
	    message: string;
	    vpcs: VPCView[];

	    static createFrom(source: any = {}) {
	        return new VPCListResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.vpcs = this.convertValues(source["vpcs"], VPCView);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class VSwitchView {
	    vswitchId: string;
	    vswitchName: string;
	    cidrBlock: string;
	    zoneId: string;
	    status: string;
	    vpcId: string;
	    creationTime: string;

	    static createFrom(source: any = {}) {
	        return new VSwitchView(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vswitchId = source["vswitchId"];
	        this.vswitchName = source["vswitchName"];
	        this.cidrBlock = source["cidrBlock"];
	        this.zoneId = source["zoneId"];
	        this.status = source["status"];
	        this.vpcId = source["vpcId"];
	        this.creationTime = source["creationTime"];
	    }
	}
	export class VSwitchListResult {
	    success: boolean;
	    message: string;
	    vswitches: VSwitchView[];

	    static createFrom(source: any = {}) {
	        return new VSwitchListResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.vswitches = this.convertValues(source["vswitches"], VSwitchView);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SLBView {
	    loadBalancerId: string;
	    loadBalancerName: string;
	    address: string;
	    addressType: string;
	    status: string;
	    regionId: string;
	    vpcId: string;
	    creationTime: string;

	    static createFrom(source: any = {}) {
	        return new SLBView(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.loadBalancerId = source["loadBalancerId"];
	        this.loadBalancerName = source["loadBalancerName"];
	        this.address = source["address"];
	        this.addressType = source["addressType"];
	        this.status = source["status"];
	        this.regionId = source["regionId"];
	        this.vpcId = source["vpcId"];
	        this.creationTime = source["creationTime"];
	    }
	}
	export class SLBDetailView {
	    loadBalancerId: string;
	    loadBalancerName: string;
	    address: string;
	    addressType: string;
	    status: string;
	    regionId: string;
	    vpcId: string;
	    vswitchId: string;
	    creationTime: string;
	    listenerCount: number;
	    bandwidth: number;

	    static createFrom(source: any = {}) {
	        return new SLBDetailView(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.loadBalancerId = source["loadBalancerId"];
	        this.loadBalancerName = source["loadBalancerName"];
	        this.address = source["address"];
	        this.addressType = source["addressType"];
	        this.status = source["status"];
	        this.regionId = source["regionId"];
	        this.vpcId = source["vpcId"];
	        this.vswitchId = source["vswitchId"];
	        this.creationTime = source["creationTime"];
	        this.listenerCount = source["listenerCount"];
	        this.bandwidth = source["bandwidth"];
	    }
	}
	export class SLBDetailResult {
	    success: boolean;
	    message: string;
	    detail?: SLBDetailView;

	    static createFrom(source: any = {}) {
	        return new SLBDetailResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.detail = this.convertValues(source["detail"], SLBDetailView);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SLBListenerView {
	    listenerPort: number;
	    listenerProtocol: string;
	    status: string;
	    bandwidth: number;

	    static createFrom(source: any = {}) {
	        return new SLBListenerView(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.listenerPort = source["listenerPort"];
	        this.listenerProtocol = source["listenerProtocol"];
	        this.status = source["status"];
	        this.bandwidth = source["bandwidth"];
	    }
	}
	export class SLBListResult {
	    success: boolean;
	    message: string;
	    slbs: SLBView[];

	    static createFrom(source: any = {}) {
	        return new SLBListResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.slbs = this.convertValues(source["slbs"], SLBView);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SLBListenerListResult {
	    success: boolean;
	    message: string;
	    listeners: SLBListenerView[];

	    static createFrom(source: any = {}) {
	        return new SLBListenerListResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.listeners = this.convertValues(source["listeners"], SLBListenerView);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

