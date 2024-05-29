<?php

namespace App\Http\Controllers;

use App\Services\GeocacheService;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class GeocacheApiController
{
    private $_service;

    public function __construct(GeocacheService $service)
    {
        $this->_service = $service;
    }

    public function all(Request $request, $isAdmin)
    {
        $search = $request->query('search', '');
        $perPage = $request->query('per_page', 10);
        $sortBy = $request->query('sort_by', '');
        $sortDirection = $request->query('sort_direction', 'desc');

        return $this->_service->all(perPage: $perPage, search: $search, sortBy: $sortBy, sortDirection: $sortDirection, isAdmin: $isAdmin);
    }

    public function allPublic (Request $request)
    {
        return response()->json($this->all($request, false), Response::HTTP_OK);
    }

    public function allAdmin (Request $request)
    {
        return response()->json($this->all($request, true), Response::HTTP_OK);
    }

    public function getByIdAdmin(int $id)
    {
        return response()->json($this->_service->getById(id: $id, isAdmin: true), Response::HTTP_OK);
    }

    public function add(Request $request)
    {
        $data = $request->all();
        $geocache = $this->_service->add($data);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["data" => $geocache];
    }

    public function update(Request $request, int $id)
    {
        $data = $request->all();
        $geocache = $this->_service->update($id, $data);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["data" => $geocache];
    }

    public function delete(int $id)
    {
        $deletedGeocache = $this->_service->delete($id);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["deleted" => $deletedGeocache];
    }
}
