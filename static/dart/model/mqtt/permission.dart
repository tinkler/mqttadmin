// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT.
import '../../http.dart';
import './const.dart';
import './role.dart' as $role show Role;
import './user.dart' as $user show User;


class Permission {
	
	int id = 0;
	
	String name = "";
	
	Future<void> savePermission(
		
	) async {
		var response = await D.instance.dio.post('$modelUrlPrefix/permission/permission/save-permission', data: {
			"data": toJson(),
			"args": {  }
		});
		if (response.data['code'] == 0) {
			var respModel = Permission.fromJson(response.data['data']['data']);
			assign(respModel);
			
		}
		
	}
	Future<void> deletePermission(
		
	) async {
		var response = await D.instance.dio.post('$modelUrlPrefix/permission/permission/delete-permission', data: {
			"data": toJson(),
			"args": {  }
		});
		if (response.data['code'] == 0) {
			var respModel = Permission.fromJson(response.data['data']['data']);
			assign(respModel);
			
		}
		
	}
	
	Permission();

	assign(Permission other) {
		
		id = other.id;
		
		name = other.name;
		
	}

	Map<String, dynamic> toJson() {
		return {
			
			"id": id,
			
			"name": name,
			
		};
	}
	Permission.fromJson(Map<String, dynamic> json) {
		
		id = json["id"];
		
		name = json["name"];
		
	}
}

class RolePermission {
	
	int id = 0;
	
	$role.Role? role;
	
	Permission? permission;
	
	Future<void> savePermission(
		
	) async {
		var response = await D.instance.dio.post('$modelUrlPrefix/permission/role_permission/save-permission', data: {
			"data": toJson(),
			"args": {  }
		});
		if (response.data['code'] == 0) {
			var respModel = RolePermission.fromJson(response.data['data']['data']);
			assign(respModel);
			
		}
		
	}
	Future<void> deletePermission(
		
	) async {
		var response = await D.instance.dio.post('$modelUrlPrefix/permission/role_permission/delete-permission', data: {
			"data": toJson(),
			"args": {  }
		});
		if (response.data['code'] == 0) {
			var respModel = RolePermission.fromJson(response.data['data']['data']);
			assign(respModel);
			
		}
		
	}
	
	RolePermission();

	assign(RolePermission other) {
		
		id = other.id;
		
		role = other.role;
		
		permission = other.permission;
		
	}

	Map<String, dynamic> toJson() {
		return {
			
			"id": id,
			
			"role": role != null ? role!.toJson() : null,
			
			"permission": permission != null ? permission!.toJson() : null,
			
		};
	}
	RolePermission.fromJson(Map<String, dynamic> json) {
		
		id = json["id"];
		
		role = json["role"] == null ? $role.Role() : $role.Role.fromJson(json["role"]);
		
		permission = json["permission"] == null ? Permission() : Permission.fromJson(json["permission"]);
		
	}
}

class UserPermission {
	
	int id = 0;
	
	$user.User? user;
	
	Permission? permission;
	
	Future<void> savePermission(
		
	) async {
		var response = await D.instance.dio.post('$modelUrlPrefix/permission/user_permission/save-permission', data: {
			"data": toJson(),
			"args": {  }
		});
		if (response.data['code'] == 0) {
			var respModel = UserPermission.fromJson(response.data['data']['data']);
			assign(respModel);
			
		}
		
	}
	Future<void> deletePermission(
		
	) async {
		var response = await D.instance.dio.post('$modelUrlPrefix/permission/user_permission/delete-permission', data: {
			"data": toJson(),
			"args": {  }
		});
		if (response.data['code'] == 0) {
			var respModel = UserPermission.fromJson(response.data['data']['data']);
			assign(respModel);
			
		}
		
	}
	
	UserPermission();

	assign(UserPermission other) {
		
		id = other.id;
		
		user = other.user;
		
		permission = other.permission;
		
	}

	Map<String, dynamic> toJson() {
		return {
			
			"id": id,
			
			"user": user != null ? user!.toJson() : null,
			
			"permission": permission != null ? permission!.toJson() : null,
			
		};
	}
	UserPermission.fromJson(Map<String, dynamic> json) {
		
		id = json["id"];
		
		user = json["user"] == null ? $user.User() : $user.User.fromJson(json["user"]);
		
		permission = json["permission"] == null ? Permission() : Permission.fromJson(json["permission"]);
		
	}
}


