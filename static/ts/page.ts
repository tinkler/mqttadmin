// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT.
import { User } from './user';


export interface Page {
	
	page: number;
	
	perPage: number;
	
	total: number;
	
	
	fetchUser(): Promise<User[]> ;
	
}

export interface PageRow {
	
	rowNo: number;
	
	chapters: ;
	
	option: ;
	
	
}

export interface Chapter {
	
	index: number;
	
	name: string;
	
	
}



export function Page(): Page {
	
	return {
		
		page: 0,
		
		perPage: 0,
		
		total: 0,
		
		
		fetchUser(): Promise<User[]>  {
			
			return postPage(this, 'fetch-user', {  }).then((res: { data: any }) => res.data as User[]);
			
		},
		
		
	};
	
}

// post data by restful api

function postPage(page: Page, method: string, args: {}): Promise<any> {
	const xhr = new XMLHttpRequest();
	xhr.open("POST", `/page//page/page/${method}`, true);
	xhr.setRequestHeader("Content-Type", "application/json");
	return new Promise((resolve, reject) => {
		xhr.onload = () => {
			if (xhr.status === 200) {
				resolve(xhr.response);
			} else {
				reject(new Error(xhr.statusText));
			}
		};
		xhr.onerror = () => {
			reject(new Error(xhr.statusText));
		};
		xhr.send(JSON.stringify({ data: page, args }));
	});
}

export function PageRow(): PageRow {
	
	return {
		
		rowNo: 0,
		
		chapters: ,
		
		option: ,
		
		
		
	};
	
}

// post data by restful api

function postPageRow(pageRow: PageRow, method: string, args: {}): Promise<any> {
	const xhr = new XMLHttpRequest();
	xhr.open("POST", `/page//page_row/page_row/${method}`, true);
	xhr.setRequestHeader("Content-Type", "application/json");
	return new Promise((resolve, reject) => {
		xhr.onload = () => {
			if (xhr.status === 200) {
				resolve(xhr.response);
			} else {
				reject(new Error(xhr.statusText));
			}
		};
		xhr.onerror = () => {
			reject(new Error(xhr.statusText));
		};
		xhr.send(JSON.stringify({ data: pageRow, args }));
	});
}

export function Chapter(): Chapter {
	
	return {
		
		index: 0,
		
		name: "",
		
		
		
	};
	
}

// post data by restful api

function postChapter(chapter: Chapter, method: string, args: {}): Promise<any> {
	const xhr = new XMLHttpRequest();
	xhr.open("POST", `/page//chapter/chapter/${method}`, true);
	xhr.setRequestHeader("Content-Type", "application/json");
	return new Promise((resolve, reject) => {
		xhr.onload = () => {
			if (xhr.status === 200) {
				resolve(xhr.response);
			} else {
				reject(new Error(xhr.statusText));
			}
		};
		xhr.onerror = () => {
			reject(new Error(xhr.statusText));
		};
		xhr.send(JSON.stringify({ data: chapter, args }));
	});
}

