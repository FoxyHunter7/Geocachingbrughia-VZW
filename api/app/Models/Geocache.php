<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Geocache extends Model
{
    use HasFactory;

    protected $hidden = [
        'id',
        'created_at',
        'updated_at'
    ];
    protected $fillable = [
        'state',
        'title',
        'geolink',
        'type',
        'difficulty',
        'terrain',
        'placed_on'
    ];
}
